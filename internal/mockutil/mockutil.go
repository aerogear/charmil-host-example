package mockutil

import (
	"context"
	"errors"
	"log"

	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"

	"github.com/aerogear/charmil-host-example/pkg/api"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
)

func NewCfgHandlerMock(cfg *config.Config) *config.CfgHandler {
	cfgHandler, err := config.NewHandler(cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfgHandler.FilePath = config.TestPath

	return cfgHandler
}

func NewConnectionMock(conn *connection.KeycloakConnection, apiClient *kafkamgmtclient.APIClient) connection.Connection {
	return &connection.ConnectionMock{
		RefreshTokensFunc: func(ctx context.Context) error {
			if conn.Token.AccessToken == "" && conn.Token.RefreshToken == "" {
				return errors.New("")
			}
			if conn.Token.RefreshToken == "expired" {
				return errors.New("")
			}

			return nil
		},
		LogoutFunc: func(ctx context.Context) error {
			if conn.Token.AccessToken == "" && conn.Token.RefreshToken == "" {
				return errors.New("")
			}
			if conn.Token.AccessToken == "expired" && conn.Token.RefreshToken == "expired" {
				return errors.New("")
			}

			conn.CfgHandler.Cfg.AccessToken = ""
			conn.CfgHandler.Cfg.RefreshToken = ""
			conn.CfgHandler.Cfg.MasAccessToken = ""
			conn.CfgHandler.Cfg.MasRefreshToken = ""

			return nil
		},
		APIFunc: func() *api.API {
			a := &api.API{
				Kafka: func() kafkamgmtclient.DefaultApi {
					return apiClient.DefaultApi
				},
			}

			return a
		},
	}
}

func NewKafkaRequestTypeMock(name string) kafkamgmtclient.KafkaRequest {
	var kafkaReq kafkamgmtclient.KafkaRequest
	kafkaReq.SetId("1")
	kafkaReq.SetName(name)

	return kafkaReq
}
