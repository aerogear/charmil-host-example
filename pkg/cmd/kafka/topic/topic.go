package topic

import (
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic/create"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic/delete"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic/describe"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic/list"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/topic/update"
)

const (
	Name      = "name"
	Operation = "operation"
)

// NewTopicCommand gives commands that manages Kafka topics.
func NewTopicCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   f.Localizer.LocalizeByID("kafka.topic.cmd.use"),
		Short: f.Localizer.LocalizeByID("kafka.topic.cmd.shortDescription"),
		Long:  f.Localizer.LocalizeByID("kafka.topic.cmd.longDescription"),
	}

	cmd.AddCommand(
		create.NewCreateTopicCommand(f),
		list.NewListTopicCommand(f),
		delete.NewDeleteTopicCommand(f),
		describe.NewDescribeTopicCommand(f),
		update.NewUpdateTopicCommand(f),
	)

	return cmd
}
