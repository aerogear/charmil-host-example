// REST API exposed via the serve command.
package registry

import (
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/registry/create"
	"github.com/aerogear/charmil-host-example/pkg/cmd/registry/delete"
	"github.com/aerogear/charmil-host-example/pkg/cmd/registry/describe"

	"github.com/aerogear/charmil-host-example/pkg/cmd/registry/list"
	"github.com/aerogear/charmil-host-example/pkg/cmd/registry/use"
	"github.com/aerogear/charmil-host-example/pkg/profile"
	"github.com/spf13/cobra"
)

func NewServiceRegistryCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "service-registry",
		// Hidden:      !profile.DevModeEnabled(),
		Annotations: profile.DevPreviewAnnotation(),
		Short:       f.Localizer.LocalizeByID("registry.cmd.shortDescription"),
		Long:        f.Localizer.LocalizeByID("registry.cmd.longDescription"),
		Example:     f.Localizer.LocalizeByID("registry.cmd.example"),
		Args:        cobra.MinimumNArgs(1),
	}

	// add sub-commands
	cmd.AddCommand(
		create.NewCreateCommand(f),
		describe.NewDescribeCommand(f),
		delete.NewDeleteCommand(f),
		list.NewListCommand(f),
		use.NewUseCommand(f),
	)

	// profile.ApplyDevPreviewLabel(cmd)

	return cmd
}
