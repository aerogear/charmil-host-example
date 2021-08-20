package root

import (
	"flag"

	"github.com/aerogear/charmil-host-example/internal/build"

	"github.com/aerogear/charmil-host-example/pkg/cmd/login"
	"github.com/aerogear/charmil-host-example/pkg/cmd/status"
	"github.com/aerogear/charmil-host-example/pkg/cmd/whoami"
	"github.com/aerogear/charmil-host-example/pkg/config"

	pluginfactory "github.com/aerogear/charmil-plugin-example/pkg/cmd/factory"
	pluginCfg "github.com/aerogear/charmil-plugin-example/pkg/config"

	"github.com/aerogear/charmil-plugin-example/pkg/cmd/registry"

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
	// cmd.AddCommand(config.NewConfigCommand(f))

	if !f.CfgHandler.Cfg.HasServiceConfigMap() {
		f.CfgHandler.Cfg.Services = &config.ServiceConfigMap{
			Kafka:           &config.KafkaConfig{},
			ServiceRegistry: &pluginCfg.Config{},
		}
	}

	// Creates a config handler instance for plugin by passing the suitable config field.
	// This line is responsible for interaction between plugin config and the host config file.
	pCfgHandler := &pluginCfg.CfgHandler{
		Cfg: f.CfgHandler.Cfg.Services.ServiceRegistry,
	}

	// Creates a plugin factory instance by passing the newly created config handler instance
	pFactory := pluginfactory.New(build.Version, nil, pCfgHandler)

	// pluginBuilder := pluginConnection.NewBuilder()

	// if cfg.AccessToken != "" {
	// 	pluginBuilder.WithAccessToken(cfg.AccessToken)
	// }
	// if cfg.RefreshToken != "" {
	// 	pluginBuilder.WithRefreshToken(cfg.RefreshToken)
	// }
	// if cfg.MasAccessToken != "" {
	// 	pluginBuilder.WithMASAccessToken(cfg.MasAccessToken)
	// }
	// if cfg.MasRefreshToken != "" {
	// 	pluginBuilder.WithMASRefreshToken(cfg.MasRefreshToken)
	// }
	// if cfg.ClientID != "" {
	// 	pluginBuilder.WithClientID(cfg.ClientID)
	// }
	// if cfg.Scopes != nil {
	// 	pluginBuilder.WithScopes(cfg.Scopes...)
	// }
	// if cfg.APIUrl != "" {
	// 	pluginBuilder.WithURL(cfg.APIUrl)
	// }
	// if cfg.AuthURL == "" {
	// 	cfg.AuthURL = build.ProductionAuthURL
	// }
	// pluginBuilder.WithAuthURL(cfg.AuthURL)

	// if cfg.MasAuthURL == "" {
	// 	cfg.MasAuthURL = build.ProductionMasAuthURL
	// }
	// pluginBuilder.WithMASAuthURL(cfg.MasAuthURL)

	// pluginBuilder.WithInsecure(cfg.Insecure)

	// pluginBuilder.WithConfig(cfgFile)

	// cmd.AddCommand(registry.NewServiceRegistryCommand(pFactory, pluginBuilder))

	// Early stage/dev preview commands
	cmd.AddCommand(registry.NewServiceRegistryCommand(pFactory))

	return cmd
}
