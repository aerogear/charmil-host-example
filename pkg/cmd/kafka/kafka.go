// Package kafka instance contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package kafka

import (
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/consumergroup"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/create"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/delete"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/describe"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/list"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/use"
)

func NewKafkaCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   f.Localizer.LocalizeByID("kafka.cmd.use"),
		Short: f.Localizer.LocalizeByID("kafka.cmd.shortDescription"),
		Args:  cobra.MinimumNArgs(1),
	}

	// add sub-commands
	cmd.AddCommand(
		create.NewCreateCommand(f),
		describe.NewDescribeCommand(f),
		delete.NewDeleteCommand(f),
		list.NewListCommand(f),
		use.NewUseCommand(f),
		topic.NewTopicCommand(f),
		consumergroup.NewConsumerGroupCommand(f),
	)

	return cmd
}
