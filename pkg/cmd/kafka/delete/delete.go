package delete

import (
	"context"
	"errors"
	"fmt"

	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/kafka"
	"github.com/aerogear/charmil/core/utils/iostreams"

	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type options struct {
	id    string
	name  string
	force bool

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

// NewDeleteCommand command for deleting kafkas.
func NewDeleteCommand(f *factory.Factory) *cobra.Command {
	opts := &options{
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.delete.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.delete.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.delete.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.delete.cmd.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.IO.CanPrompt() && !opts.force {
				return flag.RequiredWhenNonInteractiveError("yes")
			}

			if len(args) > 0 {
				opts.name = args[0]
			}

			if opts.name != "" && opts.id != "" {
				return errors.New(opts.localizer.LocalizeByID("service.error.idAndNameCannotBeUsed"))
			}

			if opts.id != "" || opts.name != "" {
				return runDelete(opts)
			}

			cfg, err := opts.Config.Load()
			if err != nil {
				return err
			}

			var kafkaConfig *config.KafkaConfig
			if cfg.Services.Kafka == kafkaConfig || cfg.Services.Kafka.ClusterID == "" {
				return errors.New(opts.localizer.LocalizeByID("kafka.common.error.noKafkaSelected"))
			}

			opts.id = cfg.Services.Kafka.ClusterID

			return runDelete(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("kafka.delete.flag.id"))
	cmd.Flags().BoolVarP(&opts.force, "yes", "y", false, opts.localizer.LocalizeByID("kafka.delete.flag.yes"))

	return cmd
}

func runDelete(opts *options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	cfg, err := opts.Config.Load()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	api := connection.API()

	var response *kafkamgmtclient.KafkaRequest
	ctx := context.Background()
	if opts.name != "" {
		response, _, err = kafka.GetKafkaByName(ctx, api.Kafka(), opts.name)
		if err != nil {
			return err
		}
	} else {
		response, _, err = kafka.GetKafkaByID(ctx, api.Kafka(), opts.id)
		if err != nil {
			return err
		}
	}

	kafkaName := response.GetName()

	if !opts.force {
		promptConfirmName := &survey.Input{
			Message: opts.localizer.LocalizeByID("kafka.delete.input.confirmName.message"),
		}

		var confirmedKafkaName string
		err = survey.AskOne(promptConfirmName, &confirmedKafkaName)
		if err != nil {
			return err
		}

		if confirmedKafkaName != kafkaName {
			logger.Info(opts.localizer.LocalizeByID("kafka.delete.log.info.incorrectNameConfirmation"))
			return nil
		}
	}

	// delete the Kafka
	logger.Infoln(opts.localizer.LocalizeByID("kafka.delete.log.debug.deletingKafka"), fmt.Sprintf("\"%s\"", kafkaName))
	a := api.Kafka().DeleteKafkaById(context.Background(), response.GetId())
	a = a.Async(true)
	_, _, err = a.Execute()

	if err != nil {
		return err
	}

	logger.Info(opts.localizer.LocalizeByID("kafka.delete.log.info.deleteSuccess", localize.NewEntry("Name", kafkaName)))

	currentKafka := cfg.Services.Kafka
	// this is not the current cluster, our work here is done
	if currentKafka == nil || currentKafka.ClusterID != response.GetId() {
		return nil
	}

	// the Kafka that was deleted is set as the user's current cluster
	// since it was deleted it should be removed from the config
	cfg.Services.Kafka = nil
	err = opts.Config.Save(cfg)
	if err != nil {
		return err
	}

	return nil
}
