package update

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/localize"

	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	topicutil "github.com/aerogear/charmil-host-example/pkg/kafka/topic"

	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"

	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	kafkainstanceclient "github.com/redhat-developer/app-services-sdk-go/kafkainstance/apiv1internal/client"
	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

var (
	partitionCount     int32
	retentionPeriodMs  int
	retentionSizeBytes int
)

type Options struct {
	topicName         string
	partitionsStr     string
	retentionMsStr    string
	retentionBytesStr string
	kafkaID           string
	outputFormat      string
	interactive       bool
	cleanupPolicy     string

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

// NewUpdateTopicCommand gets a new command for updating a kafka topic.
// nolint:funlen
func NewUpdateTopicCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Connection: f.Connection,
		CfgHandler: f.CfgHandler,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.topic.update.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.topic.update.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.topic.update.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.topic.update.cmd.example"),
		Args:    cobra.ExactValidArgs(1),
		// Dynamic completion of the topic name
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cmdutil.FilterValidTopicNameArgs(f, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			validator := topicutil.Validator{
				Localizer: opts.localizer,
			}

			if !opts.IO.CanPrompt() && opts.retentionMsStr == "" && opts.partitionsStr == "" && opts.retentionBytesStr == "" {
				return errors.New(opts.localizer.LocalizeByID("argument.error.requiredWhenNonInteractive", localize.NewEntry("Argument", "name")))
			} else if opts.retentionMsStr == "" && opts.partitionsStr == "" && opts.retentionBytesStr == "" && opts.cleanupPolicy == "" {
				opts.interactive = true
			}

			opts.topicName = args[0]

			if !opts.interactive {

				// nolint:govet
				logger, err := opts.Logger()
				if err != nil {
					return err
				}

				if opts.retentionMsStr == "" && opts.partitionsStr == "" && opts.retentionBytesStr == "" && opts.cleanupPolicy == "" {
					logger.Info(opts.localizer.LocalizeByID("kafka.topic.update.log.info.nothingToUpdate"))
					return nil
				}

				if err = validator.ValidateName(opts.topicName); err != nil {
					return err
				}

				// check that a valid --cleanup-policy flag value is used
				if opts.cleanupPolicy != "" {
					validPolicy := flagutil.IsValidInput(opts.cleanupPolicy, topicutil.ValidCleanupPolicies...)
					if !validPolicy {
						return flag.InvalidValueError("cleanup-policy", opts.cleanupPolicy, topicutil.ValidCleanupPolicies...)
					}
				}

			}

			if err = flag.ValidateOutput(opts.outputFormat); err != nil {
				return err
			}

			// check if the partition flag is set
			if opts.partitionsStr != "" {
				// nolint:govet
				partitionCount, err = topicutil.ConvertPartitionsToInt(opts.partitionsStr)
				if err != nil {
					return err
				}

				if err = validator.ValidatePartitionsN(partitionCount); err != nil {
					return err
				}
			}

			if opts.retentionMsStr != "" {
				retentionPeriodMs, err = topicutil.ConvertRetentionMsToInt(opts.retentionMsStr)
				if err != nil {
					return err
				}

				if err = validator.ValidateMessageRetentionPeriod(retentionPeriodMs); err != nil {
					return err
				}
			}

			if opts.retentionBytesStr != "" {
				retentionSizeBytes, err = topicutil.ConvertRetentionBytesToInt(opts.retentionBytesStr)
				if err != nil {
					return err
				}

				if err = validator.ValidateMessageRetentionSize(retentionSizeBytes); err != nil {
					return err
				}
			}

			if !f.CfgHandler.Cfg.HasKafka() {
				return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.noKafkaSelected"))
			}

			opts.kafkaID = opts.CfgHandler.Cfg.Services.Kafka.ClusterID

			return runCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("kafka.topic.common.flag.output.description"))
	cmd.Flags().StringVar(&opts.retentionMsStr, "retention-ms", "", opts.localizer.LocalizeByID("kafka.topic.common.input.retentionMs.description"))
	cmd.Flags().StringVar(&opts.retentionBytesStr, "retention-bytes", "", opts.localizer.LocalizeByID("kafka.topic.common.input.retentionBytes.description"))
	cmd.Flags().StringVar(&opts.cleanupPolicy, "cleanup-policy", "", opts.localizer.LocalizeByID("kafka.topic.common.input.cleanupPolicy.description"))
	cmd.Flags().StringVar(&opts.partitionsStr, "partitions", "", opts.localizer.LocalizeByID("kafka.topic.common.input.partitions.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	flagutil.EnableStaticFlagCompletion(cmd, "cleanup-policy", topicutil.ValidCleanupPolicies)

	return cmd
}

// nolint:funlen
func runCmd(opts *Options) error {
	if opts.interactive {
		// run the update command interactively
		err := runInteractivePrompt(opts)
		if err != nil {
			return err
		}

		if opts.retentionMsStr != "" {
			retentionPeriodMs, err = topicutil.ConvertRetentionMsToInt(opts.retentionMsStr)
			if err != nil {
				return err
			}
		}

		if opts.retentionBytesStr != "" {
			retentionSizeBytes, err = topicutil.ConvertRetentionBytesToInt(opts.retentionBytesStr)
			if err != nil {
				return err
			}
		}

		if opts.partitionsStr != "" {
			partitionCount, err = topicutil.ConvertPartitionsToInt(opts.partitionsStr)
			if err != nil {
				return err
			}
		}

	}

	conn, err := opts.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}
	api, kafkaInstance, err := conn.API().KafkaAdmin(opts.kafkaID)
	if err != nil {
		return err
	}

	// track if any values have changed
	var needsUpdate bool

	topic, httpRes, err := api.TopicsApi.GetTopic(context.Background(), opts.topicName).Execute()

	topicNameTmplPair := localize.NewEntry("TopicName", opts.topicName)
	kafkaNameTmplPair := localize.NewEntry("InstanceName", kafkaInstance.GetName())
	if err != nil {
		if httpRes == nil {
			return err
		}
		if httpRes.StatusCode == 404 {
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.topicNotFoundError", topicNameTmplPair, kafkaNameTmplPair))
		}
	}

	// map to store the config entries which will be updated
	configEntryMap := map[string]*string{}

	updateTopicReq := api.TopicsApi.UpdateTopic(context.Background(), opts.topicName)

	topicSettings := &kafkainstanceclient.UpdateTopicInput{}

	if opts.retentionMsStr != "" {
		needsUpdate = true
		configEntryMap[topicutil.RetentionMsKey] = &opts.retentionMsStr
	}

	if opts.retentionBytesStr != "" {
		needsUpdate = true
		configEntryMap[topicutil.RetentionSizeKey] = &opts.retentionBytesStr
	}

	if opts.cleanupPolicy != "" && strings.Compare(opts.cleanupPolicy, topicutil.GetConfigValue(topic.GetConfig(), topicutil.CleanupPolicy)) != 0 {
		needsUpdate = true
		configEntryMap[topicutil.CleanupPolicy] = &opts.cleanupPolicy
	}

	if opts.partitionsStr != "" {
		needsUpdate = true
		topicSettings.SetNumPartitions(partitionCount)
	}

	if !needsUpdate {
		logger.Info(opts.localizer.LocalizeByID("kafka.topic.update.log.info.nothingToUpdate"))
		return nil
	}

	if len(configEntryMap) > 0 {
		configEntries := topicutil.CreateConfigEntries(configEntryMap)
		topicSettings.SetConfig(*configEntries)
	}

	updateTopicReq = updateTopicReq.UpdateTopicInput(*topicSettings)

	// update the topic
	response, httpRes, err := updateTopicReq.Execute()
	// handle error
	if err != nil {
		if httpRes == nil {
			return err
		}

		operationTmplPair := localize.NewEntry("Operation", "update")
		switch httpRes.StatusCode {
		case 404:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.notFoundError", topicNameTmplPair, kafkaNameTmplPair))
		case 401:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.unauthorized", operationTmplPair))
		case 403:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.forbidden", operationTmplPair))
		case 500:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.internalServerError"))
		case 503:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.unableToConnectToKafka", localize.NewEntry("Name", kafkaInstance.GetName())))
		default:
			return err
		}
	}

	logger.Info(opts.localizer.LocalizeByID("kafka.topic.update.log.info.topicUpdated", topicNameTmplPair, kafkaNameTmplPair))

	switch opts.outputFormat {
	case dump.JSONFormat:
		data, _ := json.Marshal(response)
		_ = dump.JSON(opts.IO.Out, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(response)
		_ = dump.YAML(opts.IO.Out, data)
	}

	return nil
}

func runInteractivePrompt(opts *Options) (err error) {
	conn, err := opts.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return err
	}

	api, kafkaInstance, err := conn.API().KafkaAdmin(opts.kafkaID)
	if err != nil {
		return err
	}

	// check if topic exists
	topic, httpRes, err := api.TopicsApi.GetTopic(context.Background(), opts.topicName).
		Execute()

	topicNameTmplPair := localize.NewEntry("TopicName", opts.topicName)
	kafkaNameTmplPair := localize.NewEntry("InstanceName", kafkaInstance.GetName())
	if err != nil {
		if httpRes == nil {
			return err
		}
		if httpRes.StatusCode == 404 {
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.topicNotFoundError", topicNameTmplPair, kafkaNameTmplPair))
		}
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	validator := topicutil.Validator{
		Localizer: opts.localizer,
	}

	logger.Infoln(opts.localizer.LocalizeByID("common.log.debug.startingInteractivePrompt"))

	partitionsPrompt := &survey.Input{
		Message: opts.localizer.LocalizeByID("kafka.topic.update.input.partitions.message"),
		Help:    opts.localizer.LocalizeByID("kafka.topic.update.input.partitions.help"),
	}

	validator.CurPartitions = len(*topic.Partitions)

	err = survey.AskOne(partitionsPrompt, &opts.partitionsStr, survey.WithValidator(validator.ValidatePartitionsN))
	if err != nil {
		return err
	}

	retentionMsPrompt := &survey.Input{
		Message: opts.localizer.LocalizeByID("kafka.topic.update.input.retentionMs.message"),
		Help:    opts.localizer.LocalizeByID("kafka.topic.update.input.retentionMs.help"),
	}

	err = survey.AskOne(retentionMsPrompt, &opts.retentionMsStr, survey.WithValidator(validator.ValidateMessageRetentionPeriod))
	if err != nil {
		return err
	}

	retentionBytesPrompt := &survey.Input{
		Message: opts.localizer.LocalizeByID("kafka.topic.update.input.retentionBytes.message"),
		Help:    opts.localizer.LocalizeByID("kafka.topic.update.input.retentionBytes.help"),
	}

	err = survey.AskOne(retentionBytesPrompt, &opts.retentionBytesStr, survey.WithValidator(validator.ValidateMessageRetentionSize))
	if err != nil {
		return err
	}

	cleanupPolicyPrompt := &survey.Select{
		Message: opts.localizer.LocalizeByID("kafka.topic.update.input.cleanupPolicy.message"),
		Help:    opts.localizer.LocalizeByID("kafka.topic.update.input.cleanupPolicy.help"),
		Options: topicutil.ValidCleanupPolicies,
		Default: topicutil.GetConfigValue(topic.GetConfig(), topicutil.CleanupPolicy),
	}

	err = survey.AskOne(cleanupPolicyPrompt, &opts.cleanupPolicy)
	if err != nil {
		return err
	}

	return nil
}
