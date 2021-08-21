package factory

import (
	"context"
	"net/http"

	"github.com/aerogear/charmil-host-example/internal/build"
	"github.com/aerogear/charmil-host-example/pkg/cmd/debug"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/httputil"

	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/aerogear/charmil/core/utils/logging"
)

// New creates a new command factory
// The command factory is available to all command packages
// giving centralized access to the config and API connection

// nolint:funlen
func New(cliVersion string, localizer localize.Localizer, cfgHandler *config.CfgHandler) *Factory {
	io := iostreams.System()

	var logger logging.Logger
	var conn connection.Connection

	loggerFunc := func() (logging.Logger, error) {
		if logger != nil {
			return logger, nil
		}

		loggerBuilder := logging.NewStdLoggerBuilder()
		loggerBuilder = loggerBuilder.Streams(io.Out, io.ErrOut)

		debugEnabled := debug.Enabled()
		loggerBuilder = loggerBuilder.Debug(debugEnabled)

		logger, err := loggerBuilder.Build()
		if err != nil {
			return nil, err
		}

		return logger, nil
	}

	connectionFunc := func(connectionCfg *connection.Config) (connection.Connection, error) {
		if conn != nil {
			return conn, nil
		}

		builder := connection.NewBuilder()

		if cfgHandler.Cfg.AccessToken != "" {
			builder.WithAccessToken(cfgHandler.Cfg.AccessToken)
		}
		if cfgHandler.Cfg.RefreshToken != "" {
			builder.WithRefreshToken(cfgHandler.Cfg.RefreshToken)
		}
		if cfgHandler.Cfg.MasAccessToken != "" {
			builder.WithMASAccessToken(cfgHandler.Cfg.MasAccessToken)
		}
		if cfgHandler.Cfg.MasRefreshToken != "" {
			builder.WithMASRefreshToken(cfgHandler.Cfg.MasRefreshToken)
		}
		if cfgHandler.Cfg.ClientID != "" {
			builder.WithClientID(cfgHandler.Cfg.ClientID)
		}
		if cfgHandler.Cfg.Scopes != nil {
			builder.WithScopes(cfgHandler.Cfg.Scopes...)
		}
		if cfgHandler.Cfg.APIUrl != "" {
			builder.WithURL(cfgHandler.Cfg.APIUrl)
		}
		if cfgHandler.Cfg.AuthURL == "" {
			cfgHandler.Cfg.AuthURL = build.ProductionAuthURL
		}
		builder.WithAuthURL(cfgHandler.Cfg.AuthURL)

		if cfgHandler.Cfg.MasAuthURL == "" {
			cfgHandler.Cfg.MasAuthURL = build.ProductionMasAuthURL
		}
		builder.WithMASAuthURL(cfgHandler.Cfg.MasAuthURL)

		builder.WithInsecure(cfgHandler.Cfg.Insecure)

		builder.WithConfig(cfgHandler)

		// create a logger if it has not already been created
		logger, err := loggerFunc()
		if err != nil {
			return nil, err
		}

		transportWrapper := func(a http.RoundTripper) http.RoundTripper {
			return &httputil.LoggingRoundTripper{
				Proxied: a,
				Logger:  logger,
			}
		}

		builder.WithTransportWrapper(transportWrapper)

		builder.WithConnectionConfig(connectionCfg)

		conn, err := builder.Build()
		if err != nil {
			return nil, err
		}

		err = conn.RefreshTokens(context.TODO())
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	return &Factory{
		IOStreams:  io,
		Connection: connectionFunc,
		Logger:     loggerFunc,
		Localizer:  localizer,
		CfgHandler: cfgHandler,
	}
}
