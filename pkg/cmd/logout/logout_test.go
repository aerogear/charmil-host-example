// Package cluster contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package logout

import (
	"bytes"
	"testing"

	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/localesettings"
	"golang.org/x/text/language"

	"github.com/aerogear/charmil-host-example/internal/mockutil"

	"github.com/aerogear/charmil-host-example/pkg/config"

	"github.com/aerogear/charmil-host-example/pkg/auth/token"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"

	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/aerogear/charmil/core/utils/logging"
)

func TestNewLogoutCommand(t *testing.T) {

	locConfig := &localize.Config{
		Language: &language.English,
		Files:    localesettings.DefaultLocales,
		Format:   "toml",
	}

	localizer, _ := localize.New(locConfig)

	type args struct {
		cfg        *config.Config
		connection *connection.KeycloakConnection
	}
	tests := []struct {
		name             string
		args             args
		wantAccessToken  string
		wantRefreshToken string
	}{
		{
			name:             "Successfully logs out",
			wantAccessToken:  "",
			wantRefreshToken: "",
			args: args{
				cfg: &config.Config{
					AccessToken:  "valid",
					RefreshToken: "valid",
				},
				connection: &connection.KeycloakConnection{
					Token: &token.Token{
						AccessToken:  "valid",
						RefreshToken: "valid",
					},
				},
			},
		},
		{
			name:             "Log out is unsuccessful when tokens are expired",
			wantAccessToken:  "expired",
			wantRefreshToken: "expired",
			args: args{
				cfg: &config.Config{
					AccessToken:  "expired",
					RefreshToken: "expired",
				},
				connection: &connection.KeycloakConnection{
					Token: &token.Token{
						AccessToken:  "expired",
						RefreshToken: "expired",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.args.connection.CfgHandler = mockutil.NewCfgHandlerMock(tt.args.cfg)
		// nolint
		t.Run(tt.name, func(t *testing.T) {
			factory := &factory.Factory{
				CfgHandler: mockutil.NewCfgHandlerMock(tt.args.cfg),
				Connection: func(connectionCfg *connection.Config) (connection.Connection, error) {
					return mockutil.NewConnectionMock(tt.args.connection, nil), nil
				},
				Localizer: localizer,
				Logger: func() (logging.Logger, error) {
					loggerBuilder := logging.NewStdLoggerBuilder()
					loggerBuilder = loggerBuilder.Debug(true)
					logger, err := loggerBuilder.Build()
					if err != nil {
						return nil, err
					}

					return logger, nil
				},
			}

			cmd := NewLogoutCommand(factory)
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			_ = cmd.Execute()

			if factory.CfgHandler.Cfg.AccessToken != tt.wantAccessToken && factory.CfgHandler.Cfg.RefreshToken != tt.wantRefreshToken {
				t.Errorf("Expected access token and refresh tokens to be cleared in config")
			}
		})
	}
}
