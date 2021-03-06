package describe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	cgutil "github.com/aerogear/charmil-host-example/pkg/kafka/consumergroup"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	kafkainstanceclient "github.com/redhat-developer/app-services-sdk-go/kafkainstance/apiv1internal/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil/core/utils/color"
)

type Options struct {
	kafkaID      string
	outputFormat string
	id           string

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	localizer  localize.Localizer
}

type consumerRow struct {
	MemberID      string `json:"memberId,omitempty" header:"Member ID"`
	Partition     int    `json:"partition,omitempty" header:"Partition"`
	Topic         string `json:"topic,omitempty" header:"Topic"`
	LogEndOffset  int    `json:"logEndOffset,omitempty" header:"Log end offset"`
	CurrentOffset int    `json:"offset,omitempty" header:"Current offset"`
	OffsetLag     int    `json:"lag,omitempty" header:"Offset lag"`
}

// NewDescribeConsumerGroupCommand gets a new command for describing a consumer group.
func NewDescribeConsumerGroupCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Connection: f.Connection,
		CfgHandler: f.CfgHandler,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}
	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.consumerGroup.describe.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.consumerGroup.describe.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.consumerGroup.describe.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.consumerGroup.describe.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if opts.outputFormat != "" {
				if err = flag.ValidateOutput(opts.outputFormat); err != nil {
					return err
				}
			}

			if opts.kafkaID != "" {
				return runCmd(opts)
			}

			if !f.CfgHandler.Cfg.HasKafka() {
				return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.noKafkaSelected"))
			}

			opts.kafkaID = opts.CfgHandler.Cfg.Services.Kafka.ClusterID

			return runCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "", opts.localizer.LocalizeByID("kafka.consumerGroup.common.flag.output.description"))
	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("kafka.consumerGroup.common.flag.id.description", localize.NewEntry("Action", "view")))
	_ = cmd.MarkFlagRequired("id")

	// flag based completions for ID
	_ = cmd.RegisterFlagCompletionFunc("id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cmdutil.FilterValidConsumerGroupIDs(f, toComplete)
	})

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

	ctx := context.Background()

	consumerGroupData, httpRes, err := api.GroupsApi.GetConsumerGroupById(ctx, opts.id).Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		cgIDPair := localize.NewEntry("ID", opts.id)
		kafkaNameTmplPair := localize.NewEntry("InstanceName", kafkaInstance.GetName())
		operationTmplPair := localize.NewEntry("Operation", "view")

		switch httpRes.StatusCode {
		case 404:
			return errors.New(opts.localizer.LocalizeByID("kafka.consumerGroup.common.error.notFoundError", cgIDPair, kafkaNameTmplPair))
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

	stdout := opts.IO.Out
	switch opts.outputFormat {
	case dump.JSONFormat:
		data, _ := json.Marshal(consumerGroupData)
		_ = dump.JSON(stdout, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(consumerGroupData)
		_ = dump.YAML(stdout, data)
	default:
		printConsumerGroupDetails(stdout, consumerGroupData, opts.localizer)
	}

	return nil
}

func mapConsumerGroupDescribeToTableFormat(consumers []kafkainstanceclient.Consumer) []consumerRow {
	rows := []consumerRow{}

	for _, consumer := range consumers {

		row := consumerRow{
			Partition:     int(consumer.GetPartition()),
			Topic:         consumer.GetTopic(),
			MemberID:      consumer.GetMemberId(),
			LogEndOffset:  int(consumer.GetLogEndOffset()),
			CurrentOffset: int(consumer.GetOffset()),
			OffsetLag:     int(consumer.GetLag()),
		}

		if consumer.GetMemberId() == "" {
			row.MemberID = color.Bold("unconsumed")
		}

		rows = append(rows, row)
	}

	// sort members by partition number
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Partition < rows[j].Partition
	})

	return rows
}

// print the consumer grooup details
func printConsumerGroupDetails(w io.Writer, consumerGroupData kafkainstanceclient.ConsumerGroup, localizer localize.Localizer) {
	fmt.Fprintln(w, "")
	consumers := consumerGroupData.GetConsumers()

	activeMembersCount := cgutil.GetActiveConsumersCount(consumers)
	partitionsWithLagCount := cgutil.GetPartitionsWithLag(consumers)
	unconsumedPartitions := cgutil.GetUnconsumedPartitions(consumers)

	fmt.Fprintln(w, color.Bold(localizer.LocalizeByID("kafka.consumerGroup.describe.output.activeMembers")), activeMembersCount, "\t", color.Bold(localizer.LocalizeByID("kafka.consumerGroup.describe.output.partitionsWithLag")), partitionsWithLagCount, "\t", color.Bold(localizer.LocalizeByID("kafka.consumerGroup.describe.output.unconsumedPartitions")), unconsumedPartitions)
	fmt.Fprintln(w, "")

	rows := mapConsumerGroupDescribeToTableFormat(consumers)
	dump.Table(w, rows)
}
