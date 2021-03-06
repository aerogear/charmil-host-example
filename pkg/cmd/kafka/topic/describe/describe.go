package describe

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"

	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"

	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	topicName    string
	kafkaID      string
	outputFormat string

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

// NewDescribeTopicCommand gets a new command for describing a kafka topic.
func NewDescribeTopicCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Connection: f.Connection,
		CfgHandler: f.CfgHandler,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.topic.describe.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.topic.describe.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.topic.describe.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.topic.describe.cmd.example"),
		Args:    cobra.ExactValidArgs(1),
		// dynamic completion of topic names
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cmdutil.FilterValidTopicNameArgs(f, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				opts.topicName = args[0]
			}

			if opts.outputFormat != "" {
				if err = flag.ValidateOutput(opts.outputFormat); err != nil {
					return err
				}
			}

			if opts.kafkaID != "" {
				return runCmd(opts)
			}

			if !f.CfgHandler.Cfg.HasKafka() {
				return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.noKafkaSelected"))
			}

			opts.kafkaID = opts.CfgHandler.Cfg.Services.Kafka.ClusterID

			return runCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("kafka.topic.common.flag.output.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runCmd(opts *Options) error {
	conn, err := opts.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return err
	}

	api, kafkaInstance, err := conn.API().KafkaAdmin(opts.kafkaID)
	if err != nil {
		return err
	}

	// fetch the topic
	topicResponse, httpRes, err := api.TopicsApi.
		GetTopic(context.Background(), opts.topicName).
		Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		topicNameTmplPair := localize.NewEntry("TopicName", opts.topicName)
		kafkaNameTmplPair := localize.NewEntry("InstanceName", kafkaInstance.GetName())
		operationTmplPair := localize.NewEntry("Operation", "delete")

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

	switch opts.outputFormat {
	case dump.JSONFormat:
		data, _ := json.Marshal(topicResponse)
		_ = dump.JSON(opts.IO.Out, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(topicResponse)
		_ = dump.YAML(opts.IO.Out, data)
	}

	return nil
}
