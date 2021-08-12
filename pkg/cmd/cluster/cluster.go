package cluster

import (
	"github.com/aerogear/charmil-host-example/pkg/cmd/cluster/bind"
	"github.com/aerogear/charmil-host-example/pkg/cmd/cluster/connect"
	"github.com/aerogear/charmil-host-example/pkg/cmd/cluster/status"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/spf13/cobra"
)

// NewServiceAccountCommand creates a new command sub-group to manage service accounts
func NewClusterCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     f.Localizer.MustLocalize("cluster.cmd.use"),
		Short:   f.Localizer.MustLocalize("cluster.cmd.shortDescription"),
		Example: f.Localizer.MustLocalize("cluster.cmd.example"),
		Args:    cobra.ExactArgs(1),
	}

	cmd.AddCommand(
		status.NewStatusCommand(f),
		connect.NewConnectCommand(f),
		bind.NewBindCommand(f),
	)

	return cmd
}
