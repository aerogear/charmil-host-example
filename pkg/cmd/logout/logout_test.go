// Package cluster contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package logout

import (
	"bytes"
	"testing"

	"github.com/aerogear/charmil-host-example/pkg/connection"

	"github.com/aerogear/charmil-host-example/internal/mockutil"

	"github.com/aerogear/charmil-host-example/internal/config"

	"github.com/aerogear/charmil-host-example/pkg/auth/token"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"

	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/aerogear/charmil/core/utils/logging"
)

func TestNewLogoutCommand(t *testing.T) {
	localizer, _ := localize.New(nil)

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
		tt.args.connection.Config = mockutil.NewConfigMock(tt.args.cfg)
		// nolint
		t.Run(tt.name, func(t *testing.T) {
			factory := &factory.Factory{
				Config: mockutil.NewConfigMock(tt.args.cfg),
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

			cfg, _ := factory.Config.Load()
			if cfg.AccessToken != tt.wantAccessToken && cfg.RefreshToken != tt.wantRefreshToken {
				t.Errorf("Expected access token and refresh tokens to be cleared in config")
			}
		})
	}
}
