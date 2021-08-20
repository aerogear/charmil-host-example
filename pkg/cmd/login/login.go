// Package cluster contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package login

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"

	"github.com/aerogear/charmil-host-example/internal/build"
	"golang.org/x/oauth2"

	"github.com/aerogear/charmil-host-example/pkg/auth/login"
	"github.com/aerogear/charmil-host-example/pkg/auth/token"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/debug"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil/core/utils/iostreams"

	"github.com/aerogear/charmil-host-example/pkg/connection"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

// When the value of the `--api-gateway` option is one of the keys of this map it will be replaced by the
// corresponding value.
var apiGatewayAliases = map[string]string{
	"production": build.ProductionAPIURL,
	"prod":       build.ProductionAPIURL,
	"prd":        build.ProductionAPIURL,
	"staging":    build.StagingAPIURL,
	"stage":      build.StagingAPIURL,
	"stg":        build.StagingAPIURL,
}

// When the value of the `--auth-url` option is one of the keys of this map it will be replaced by the
// corresponding value.
var authURLAliases = map[string]string{
	"production": build.ProductionAuthURL,
	"prod":       build.ProductionAuthURL,
	"prd":        build.ProductionAuthURL,
	"staging":    build.ProductionAuthURL,
	"stage":      build.ProductionAuthURL,
	"stg":        build.ProductionAuthURL,
}

// When the value of the `--mas-auth-url` option is one of the keys of this map it will be replaced by the
// corresponding value.
var masAuthURLAliases = map[string]string{
	"production": build.ProductionMasAuthURL,
	"prod":       build.ProductionMasAuthURL,
	"prd":        build.ProductionMasAuthURL,
	"staging":    build.StagingMasAuthURL,
	"stage":      build.StagingMasAuthURL,
	"stg":        build.StagingMasAuthURL,
}

type Options struct {
	CfgHandler *config.CfgHandler
	Logger     func() (logging.Logger, error)
	Connection factory.ConnectionFunc
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	url                   string
	authURL               string
	masAuthURL            string
	clientID              string
	scopes                []string
	insecureSkipTLSVerify bool
	printURL              bool
	offlineToken          string
}

// NewLoginCmd gets the command that's log the user in
func NewLoginCmd(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("login.cmd.use"),
		Short:   opts.localizer.LocalizeByID("login.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("login.cmd.longDescription", localize.NewEntry("OfflineTokenURL", build.OfflineTokenURL)),
		Example: opts.localizer.LocalizeByID("login.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.offlineToken != "" && opts.clientID == build.DefaultClientID {
				opts.clientID = build.DefaultOfflineTokenClientID
			}

			logger, err := opts.Logger()
			if err != nil {
				return err
			}

			if opts.IO.IsSSHSession() && opts.offlineToken == "" {
				logger.Info(opts.localizer.LocalizeByID("login.log.info.sshLoginDetected", localize.NewEntry("OfflineTokenURL", build.OfflineTokenURL)))
			}

			return runLogin(opts)
		},
	}

	cmd.Flags().StringVar(&opts.url, "api-gateway", build.ProductionAPIURL, opts.localizer.LocalizeByID("login.flag.apiGateway"))
	cmd.Flags().BoolVar(&opts.insecureSkipTLSVerify, "insecure", false, opts.localizer.LocalizeByID("login.flag.insecure"))
	cmd.Flags().StringVar(&opts.clientID, "client-id", build.DefaultClientID, opts.localizer.LocalizeByID("login.flag.clientId"))
	cmd.Flags().StringVar(&opts.authURL, "auth-url", build.ProductionAuthURL, opts.localizer.LocalizeByID("login.flag.authUrl"))
	cmd.Flags().StringVar(&opts.masAuthURL, "mas-auth-url", build.ProductionMasAuthURL, opts.localizer.LocalizeByID("login.flag.masAuthUrl"))
	cmd.Flags().BoolVar(&opts.printURL, "print-sso-url", false, opts.localizer.LocalizeByID("login.flag.printSsoUrl"))
	cmd.Flags().StringArrayVar(&opts.scopes, "scope", connection.DefaultScopes, opts.localizer.LocalizeByID("login.flag.scope"))
	cmd.Flags().StringVarP(&opts.offlineToken, "token", "t", "", opts.localizer.LocalizeByID("login.flag.token", localize.NewEntry("OfflineTokenURL", build.OfflineTokenURL)))

	return cmd
}

// nolint:funlen
func runLogin(opts *Options) (err error) {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	gatewayURL, err := getURLFromAlias(opts.url, apiGatewayAliases, opts.localizer)
	if err != nil {
		return err
	}

	authURL, err := getURLFromAlias(opts.authURL, authURLAliases, opts.localizer)
	if err != nil {
		return err
	}
	opts.authURL = authURL.String()

	masAuthURL, err := getURLFromAlias(opts.masAuthURL, masAuthURLAliases, opts.localizer)
	if err != nil {
		return err
	}
	opts.masAuthURL = masAuthURL.String()

	if opts.offlineToken == "" {
		tr := createTransport(opts.insecureSkipTLSVerify)
		httpClient := oauth2.NewClient(context.Background(), nil)
		httpClient.Transport = tr

		loginExec := &login.AuthorizationCodeGrant{
			HTTPClient: httpClient,
			Scopes:     opts.scopes,
			Logger:     logger,
			IO:         opts.IO,
			CfgHandler: opts.CfgHandler,
			ClientID:   opts.clientID,
			PrintURL:   opts.printURL,
			Localizer:  opts.localizer,
		}

		ssoCfg := &login.SSOConfig{
			AuthURL:      authURL,
			RedirectPath: "sso-redhat-callback",
		}

		masSsoCfg := &login.SSOConfig{
			AuthURL:      masAuthURL,
			RedirectPath: "mas-sso-callback",
		}

		if err = loginExec.Execute(context.Background(), ssoCfg, masSsoCfg); err != nil {
			return err
		}
	}

	if opts.offlineToken != "" {
		if err = loginWithOfflineToken(opts); err != nil {
			return err
		}
	}

	opts.CfgHandler.Cfg.APIUrl = gatewayURL.String()
	opts.CfgHandler.Cfg.Insecure = opts.insecureSkipTLSVerify
	opts.CfgHandler.Cfg.ClientID = opts.clientID
	opts.CfgHandler.Cfg.AuthURL = opts.authURL
	opts.CfgHandler.Cfg.MasAuthURL = opts.masAuthURL
	opts.CfgHandler.Cfg.Scopes = opts.scopes

	username, ok := token.GetUsername(opts.CfgHandler.Cfg.AccessToken)
	logger.Info("")

	if !ok {
		logger.Info(opts.localizer.LocalizeByID("login.log.info.loginSuccessNoUsername"))
	} else {
		opts.localizer.LocalizeByID("login.log.info.loginSuccess", localize.NewEntry("Username", username))
	}

	// debug mode checks this for a version update also.
	// so we check if is enabled first so as not to print it twice
	if !debug.Enabled() {
		build.CheckForUpdate(context.Background(), logger, opts.localizer)
	}

	return nil
}

func loginWithOfflineToken(opts *Options) (err error) {

	opts.CfgHandler.Cfg.Insecure = opts.insecureSkipTLSVerify
	opts.CfgHandler.Cfg.ClientID = opts.clientID
	opts.CfgHandler.Cfg.AuthURL = opts.authURL
	opts.CfgHandler.Cfg.MasAuthURL = opts.masAuthURL
	opts.CfgHandler.Cfg.Scopes = opts.scopes
	opts.CfgHandler.Cfg.RefreshToken = opts.offlineToken
	// remove MAS-SSO tokens, as this does not support token login
	opts.CfgHandler.Cfg.MasAccessToken = ""
	opts.CfgHandler.Cfg.MasRefreshToken = ""

	_, err = opts.Connection(connection.DefaultConfigSkipMasAuth)
	return err
}

func createTransport(insecure bool) *http.Transport {
	// #nosec 402
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
}

func getURLFromAlias(urlOrAlias string, urlAliasMap map[string]string, localizer localize.Localizer) (u *url.URL, err error) {
	// If the URL value is any of the aliases then replace it with the corresponding
	// real URL:
	unparsedGatewayURL, ok := urlAliasMap[urlOrAlias]
	if !ok {
		unparsedGatewayURL = urlOrAlias
	}

	gatewayURL, err := url.ParseRequestURI(unparsedGatewayURL)
	if err != nil {
		return nil, err
	}
	if gatewayURL.Scheme != "http" && gatewayURL.Scheme != "https" {
		err = errors.New(localizer.LocalizeByID("login.error.schemeMissingFromUrl", localize.NewEntry("URL", gatewayURL.String())))
		return nil, err
	}

	return gatewayURL, nil
}
