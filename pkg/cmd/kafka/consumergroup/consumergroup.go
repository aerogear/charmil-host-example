package consumergroup

import (
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/consumergroup/delete"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/consumergroup/describe"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/consumergroup/list"
	"github.com/spf13/cobra"
)

// NewConsumerGroupCommand creates a new command sub-group for consumer group operations
func NewConsumerGroupCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   f.Localizer.MustLocalize("kafka.consumerGroup.cmd.use"),
		Short: f.Localizer.MustLocalize("kafka.consumerGroup.cmd.shortDescription"),
		Long:  f.Localizer.MustLocalize("kafka.consumerGroup.cmd.longDescription"),
		Args:  cobra.ExactArgs(1),
	}

	cmd.AddCommand(
		list.NewListConsumerGroupCommand(f),
		delete.NewDeleteConsumerGroupCommand(f),
		describe.NewDescribeConsumerGroupCommand(f),
	)

	return cmd
}
