package list

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	topicutil "github.com/aerogear/charmil-host-example/pkg/kafka/topic"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/aerogear/charmil-host-example/pkg/connection"

	kafkainstanceclient "github.com/redhat-developer/app-services-sdk-go/kafkainstance/apiv1internal/client"

	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"

	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	CfgHandler *config.CfgHandler
	IO         *iostreams.IOStreams
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer

	kafkaID string
	output  string
	search  string
	page    int32
	size    int32
}

type topicRow struct {
	Name            string `json:"name,omitempty" header:"Name"`
	PartitionsCount int    `json:"partitions_count,omitempty" header:"Partitions"`
	RetentionTime   string `json:"retention.ms,omitempty" header:"Retention time (ms)"`
	RetentionSize   string `json:"retention.bytes,omitempty" header:"Retention size (bytes)"`
}

// NewListTopicCommand gets a new command for getting kafkas.
func NewListTopicCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.topic.list.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.topic.list.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.topic.list.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.topic.list.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.output != "" {
				if err := flag.ValidateOutput(opts.output); err != nil {
					return err
				}
			}

			if opts.page < 1 {
				return errors.New(opts.localizer.LocalizeByID("kafka.common.page.error.invalid.minValue", localize.NewEntry("Page", opts.page)))
			}

			if opts.size < 1 {
				return errors.New(opts.localizer.LocalizeByID("kafka.common.size.error.invalid.minValue", localize.NewEntry("Size", opts.size)))
			}

			if opts.search != "" {
				validator := topicutil.Validator{
					Localizer: opts.localizer,
				}
				if err := validator.ValidateSearchInput(opts.search); err != nil {
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

	cmd.Flags().StringVarP(&opts.output, "output", "o", "", opts.localizer.LocalizeByID("kafka.topic.list.flag.output.description"))
	cmd.Flags().StringVarP(&opts.search, "search", "", "", opts.localizer.LocalizeByID("kafka.topic.list.flag.search.description"))
	cmd.Flags().Int32VarP(&opts.page, "page", "", int32(cmdutil.DefaultPageNumber), opts.localizer.LocalizeByID("kafka.topic.list.flag.page.description"))
	cmd.Flags().Int32VarP(&opts.size, "size", "", int32(cmdutil.DefaultPageSize), opts.localizer.LocalizeByID("kafka.topic.list.flag.size.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runCmd(opts *Options) error {
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

	a := api.TopicsApi.GetTopics(context.Background())

	if opts.search != "" {
		logger.Infoln(opts.localizer.LocalizeByID("kafka.topic.list.log.debug.filteringTopicList", localize.NewEntry("Search", opts.search)))
		a = a.Filter(opts.search)
	}

	a = a.Size(opts.size)

	a = a.Page(opts.page)

	topicData, httpRes, err := a.Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		operationTemplatePair := localize.NewEntry("Operation", "list")

		switch httpRes.StatusCode {
		case http.StatusUnauthorized:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.list.error.unauthorized", operationTemplatePair))
		case http.StatusForbidden:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.list.error.forbidden", operationTemplatePair))
		case http.StatusInternalServerError:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.internalServerError"))
		case http.StatusServiceUnavailable:
			return errors.New(opts.localizer.LocalizeByID("kafka.topic.common.error.unableToConnectToKafka", localize.NewEntry("Name", kafkaInstance.GetName())))
		default:
			return err
		}
	}

	defer httpRes.Body.Close()

	if topicData.GetTotal() == 0 && opts.output == "" {
		logger.Info(opts.localizer.LocalizeByID("kafka.topic.list.log.info.noTopics", localize.NewEntry("InstanceName", kafkaInstance.GetName())))

		return nil
	}

	stdout := opts.IO.Out
	switch opts.output {
	case dump.JSONFormat:
		data, _ := json.Marshal(topicData)
		_ = dump.JSON(stdout, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(topicData)
		_ = dump.YAML(stdout, data)
	default:
		topics := topicData.GetItems()
		rows := mapTopicResultsToTableFormat(topics)
		dump.Table(stdout, rows)
	}

	return nil
}

func mapTopicResultsToTableFormat(topics []kafkainstanceclient.Topic) []topicRow {
	rows := []topicRow{}

	for _, t := range topics {

		row := topicRow{
			Name:            t.GetName(),
			PartitionsCount: len(t.GetPartitions()),
		}
		for _, config := range t.GetConfig() {
			unlimitedVal := "-1 (Unlimited)"

			if *config.Key == topicutil.RetentionMsKey {
				val := config.GetValue()
				if val == "-1" {
					row.RetentionTime = unlimitedVal
				} else {
					row.RetentionTime = val
				}
			}
			if *config.Key == topicutil.RetentionSizeKey {
				val := config.GetValue()
				if val == "-1" {
					row.RetentionSize = unlimitedVal
				} else {
					row.RetentionSize = val
				}
			}
		}

		rows = append(rows, row)
	}

	return rows
}
