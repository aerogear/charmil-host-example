package root

import (
	"context"
	"flag"
	"net/http"

	"github.com/aerogear/charmil-host-example/internal/build"

	"github.com/aerogear/charmil-host-example/pkg/cmd/login"
	"github.com/aerogear/charmil-host-example/pkg/cmd/status"
	"github.com/aerogear/charmil-host-example/pkg/cmd/whoami"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/httputil"

	pluginfactory "github.com/aerogear/charmil-plugin-example/pkg/cmd/factory"
	pluginCfg "github.com/aerogear/charmil-plugin-example/pkg/config"

	pluginConnection "github.com/aerogear/charmil-plugin-example/pkg/connection"

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

	// initilize plugin config in pFactory
	initPluginConfig(f, pFactory, pCfgHandler)

	// add service-registry command from plugin
	cmd.AddCommand(registry.NewServiceRegistryCommand(pFactory))

	return cmd
}

// initPluginConfig initializes configuration for plugin from host
func initPluginConfig(f *factory.Factory, pFactory *pluginfactory.Factory, pluginCfgHandler *pluginCfg.CfgHandler) {
	pluginBuilder := pluginConnection.NewBuilder()
	if f.CfgHandler.Cfg.AccessToken != "" {
		pluginBuilder.WithAccessToken(f.CfgHandler.Cfg.AccessToken)
	}
	if f.CfgHandler.Cfg.RefreshToken != "" {
		pluginBuilder.WithRefreshToken(f.CfgHandler.Cfg.RefreshToken)
	}
	if f.CfgHandler.Cfg.MasAccessToken != "" {
		pluginBuilder.WithMASAccessToken(f.CfgHandler.Cfg.MasAccessToken)
	}
	if f.CfgHandler.Cfg.MasRefreshToken != "" {
		pluginBuilder.WithMASRefreshToken(f.CfgHandler.Cfg.MasRefreshToken)
	}
	if f.CfgHandler.Cfg.ClientID != "" {
		pluginBuilder.WithClientID(f.CfgHandler.Cfg.ClientID)
	}
	if f.CfgHandler.Cfg.Scopes != nil {
		pluginBuilder.WithScopes(f.CfgHandler.Cfg.Scopes...)
	}
	if f.CfgHandler.Cfg.APIUrl != "" {
		pluginBuilder.WithURL(f.CfgHandler.Cfg.APIUrl)
	}
	if f.CfgHandler.Cfg.AuthURL == "" {
		f.CfgHandler.Cfg.AuthURL = build.ProductionAuthURL
	}
	pluginBuilder.WithAuthURL(f.CfgHandler.Cfg.AuthURL)
	if f.CfgHandler.Cfg.MasAuthURL == "" {
		f.CfgHandler.Cfg.MasAuthURL = build.ProductionMasAuthURL
	}
	pluginBuilder.WithMASAuthURL(f.CfgHandler.Cfg.MasAuthURL)
	pluginBuilder.WithInsecure(f.CfgHandler.Cfg.Insecure)
	pluginBuilder.WithConfig(pluginCfgHandler)

	pluginConnectionFunction := func(connectionCfg *pluginConnection.Config) (pluginConnection.Connection, error) {
		transportWrapper := func(a http.RoundTripper) http.RoundTripper {
			return &httputil.LoggingRoundTripper{
				Proxied: a,
			}
		}

		pluginBuilder.WithTransportWrapper(transportWrapper)

		pluginBuilder.WithConnectionConfig(connectionCfg)

		conn, err := pluginBuilder.Build()
		if err != nil {
			return nil, err
		}

		err = conn.RefreshTokens(context.TODO())
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	pFactory.Connection = pluginConnectionFunction
}
