package list

import (
	"context"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"

	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil-host-example/pkg/kafka/consumergroup"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	kafkainstanceclient "github.com/redhat-developer/app-services-sdk-go/kafkainstance/apiv1internal/client"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	output  string
	kafkaID string
	topic   string
	search  string
	page    int32
	size    int32
}

type consumerGroupRow struct {
	ConsumerGroupID   string `json:"groupId,omitempty" header:"Consumer group ID"`
	ActiveMembers     int    `json:"active_members,omitempty" header:"Active members"`
	PartitionsWithLag int    `json:"lag,omitempty" header:"Partitions with lag"`
}

// NewListConsumerGroupCommand creates a new command to list consumer groups
func NewListConsumerGroupCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.consumerGroup.list.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.consumerGroup.list.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.consumerGroup.list.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.consumerGroup.list.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.output != "" && !flagutil.IsValidInput(opts.output, flagutil.ValidOutputFormats...) {
				return flag.InvalidValueError("output", opts.output, flagutil.ValidOutputFormats...)
			}

			if opts.page < 1 {
				return errors.New(opts.localizer.LocalizeByID("kafka.common.validation.page.error.invalid.minValue", localize.NewEntry("Page", opts.page)))
			}

			if opts.size < 1 {
				return errors.New(opts.localizer.LocalizeByID("kafka.common.validation.size.error.invalid.minValue", localize.NewEntry("Size", opts.size)))
			}

			cfg, err := opts.Config.Load()
			if err != nil {
				return err
			}

			if !cfg.HasKafka() {
				return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.noKafkaSelected"))
			}

			opts.kafkaID = cfg.Services.Kafka.ClusterID

			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", "", opts.localizer.LocalizeByID("kafka.consumerGroup.list.flag.output.description"))
	cmd.Flags().StringVar(&opts.topic, "topic", "", opts.localizer.LocalizeByID("kafka.consumerGroup.list.flag.topic.description"))
	cmd.Flags().StringVar(&opts.search, "search", "", opts.localizer.LocalizeByID("kafka.consumerGroup.list.flag.search"))
	cmd.Flags().Int32VarP(&opts.page, "page", "", int32(cmdutil.DefaultPageNumber), opts.localizer.LocalizeByID("kafka.consumerGroup.list.flag.page"))
	cmd.Flags().Int32VarP(&opts.size, "size", "", int32(cmdutil.DefaultPageSize), opts.localizer.LocalizeByID("kafka.consumerGroup.list.flag.size"))

	_ = cmd.RegisterFlagCompletionFunc("topic", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cmdutil.FilterValidTopicNameArgs(f, toComplete)
	})

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

// nolint:funlen
func runList(opts *Options) (err error) {
	conn, err := opts.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	ctx := context.Background()

	api, kafkaInstance, err := conn.API().KafkaAdmin(opts.kafkaID)
	if err != nil {
		return err
	}

	req := api.GroupsApi.GetConsumerGroups(ctx)

	if opts.topic != "" {
		req = req.Topic(opts.topic)
	}
	if opts.search != "" {
		req = req.GroupIdFilter(opts.search)
	}

	req = req.Size(opts.size)

	req = req.Page(opts.page)

	consumerGroupData, httpRes, err := req.Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		operationTmplPair := localize.NewEntry("Operation", "list")

		switch httpRes.StatusCode {
		case 401:
			return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.unauthorized", operationTmplPair))
		case 403:
			return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.forbidden", operationTmplPair))
		case 500:
			return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.internalServerError"))
		case 503:
			return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.unableToConnectToKafka", localize.NewEntry("Name", kafkaInstance.GetName())))
		default:
			return err
		}
	}

	ok, err := checkForConsumerGroups(int(consumerGroupData.GetTotal()), opts, kafkaInstance.GetName())
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	switch opts.output {
	case dump.JSONFormat:
		data, _ := json.Marshal(consumerGroupData)
		_ = dump.JSON(opts.IO.Out, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(consumerGroupData)
		_ = dump.YAML(opts.IO.Out, data)
	default:
		logger.Info("")
		consumerGroups := consumerGroupData.GetItems()
		rows := mapConsumerGroupResultsToTableFormat(consumerGroups)
		dump.Table(opts.IO.Out, rows)

		return nil
	}

	return nil
}

func mapConsumerGroupResultsToTableFormat(consumerGroups []kafkainstanceclient.ConsumerGroup) []consumerGroupRow {
	rows := []consumerGroupRow{}

	for _, t := range consumerGroups {
		consumers := t.GetConsumers()
		row := consumerGroupRow{
			ConsumerGroupID:   t.GetGroupId(),
			ActiveMembers:     consumergroup.GetActiveConsumersCount(consumers),
			PartitionsWithLag: consumergroup.GetPartitionsWithLag(consumers),
		}
		rows = append(rows, row)
	}

	return rows
}

// checks if there are any consumer groups available
// prints to stderr if not
func checkForConsumerGroups(count int, opts *Options, kafkaName string) (hasCount bool, err error) {
	logger, err := opts.Logger()
	if err != nil {
		return false, err
	}
	kafkaNameTmplPair := localize.NewEntry("InstanceName", kafkaName)
	if count == 0 && opts.output == "" {
		if opts.topic == "" {
			logger.Info(opts.localizer.LocalizeByID("kafka.consumerGroup.list.log.info.noConsumerGroups", kafkaNameTmplPair))
		} else {
			logger.Info(opts.localizer.LocalizeByID("kafka.consumerGroup.list.log.info.noConsumerGroupsForTopic", kafkaNameTmplPair, localize.NewEntry("TopicName", opts.topic)))
		}

		return false, nil
	}

	return true, nil
}
