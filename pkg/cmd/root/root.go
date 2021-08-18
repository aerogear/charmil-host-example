package root

import (
	"flag"
	"fmt"
	"os"

	"github.com/aerogear/charmil-host-example/internal/build"
	"github.com/aerogear/charmil-host-example/pkg/cmd/config"

	"github.com/aerogear/charmil-host-example/pkg/cmd/login"
	"github.com/aerogear/charmil-host-example/pkg/cmd/status"
	"github.com/aerogear/charmil-host-example/pkg/cmd/whoami"

	pluginfactory "github.com/aerogear/charmil-plugin-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-plugin-example/pkg/cmd/registry"
	pluginloc "github.com/aerogear/charmil-plugin-example/pkg/localize/goi18n"

	"github.com/aerogear/charmil-host-example/pkg/arguments"
	"github.com/aerogear/charmil-host-example/pkg/cmd/cluster"
	"github.com/aerogear/charmil-host-example/pkg/cmd/completion"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka"
	"github.com/aerogear/charmil-host-example/pkg/cmd/logout"
	"github.com/aerogear/charmil-host-example/pkg/cmd/serviceaccount"
	cliversion "github.com/aerogear/charmil-host-example/pkg/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRootCommand(f *factory.Factory, version string) *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           f.Localizer.LocalizeByID("root.cmd.use"),
		Short:         f.Localizer.LocalizeByID("root.cmd.shortDescription"),
		Long:          f.Localizer.LocalizeByID("root.cmd.longDescription"),
		Example:       f.Localizer.LocalizeByID("root.cmd.example"),
	}
	fs := cmd.PersistentFlags()
	arguments.AddDebugFlag(fs)
	// this flag comes out of the box, but has its own basic usage text, so this overrides that
	var help bool

	fs.BoolVarP(&help, "help", "h", false, f.Localizer.LocalizeByID("root.cmd.flag.help.description"))
	fs.Bool("version", false, f.Localizer.LocalizeByID("root.cmd.flag.version.description"))

	cmd.Version = version

	// cmd.SetVersionTemplate(f.Localizer.LocalizeByID("version.cmd.outputText", localize.NewEntry("Version", build.Version)))
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Child commands
	cmd.AddCommand(login.NewLoginCmd(f))
	cmd.AddCommand(logout.NewLogoutCommand(f))
	cmd.AddCommand(kafka.NewKafkaCommand(f))
	cmd.AddCommand(serviceaccount.NewServiceAccountCommand(f))
	cmd.AddCommand(cluster.NewClusterCommand(f))
	cmd.AddCommand(status.NewStatusCommand(f))
	cmd.AddCommand(completion.NewCompletionCommand(f))
	cmd.AddCommand(whoami.NewWhoAmICmd(f))
	cmd.AddCommand(cliversion.NewVersionCmd(f))
	cmd.AddCommand(config.NewConfigCommand(f))

	loc, err := pluginloc.New(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pfact := pluginfactory.New(build.Version, loc)
	cmd.AddCommand(registry.NewServiceRegistryCommand(pfact))

	// Early stage/dev preview commands
	// cmd.AddCommand(registry.NewServiceRegistryCommand(f))

	return cmd
}
