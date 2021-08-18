package use

import (
	"context"
	"errors"
	"fmt"

	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"

	"github.com/aerogear/charmil-host-example/pkg/cmdutil"

	"github.com/aerogear/charmil-host-example/pkg/kafka"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	id          string
	name        string
	interactive bool

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewUseCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.use.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.use.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.use.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.use.cmd.example"),
		Args:    cobra.RangeArgs(0, 1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cmdutil.FilterValidKafkas(f, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.name = args[0]
			} else if opts.id == "" {
				if !opts.IO.CanPrompt() {
					return errors.New(opts.localizer.LocalizeByID("kafka.use.error.idOrNameRequired"))
				}
				opts.interactive = true
			}

			if opts.name != "" && opts.id != "" {
				return errors.New(opts.localizer.LocalizeByID("service.error.idAndNameCannotBeUsed"))
			}

			return runUse(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("kafka.use.flag.id"))

	return cmd
}

func runUse(opts *Options) error {
	if opts.interactive {
		// run the use command interactively
		err := runInteractivePrompt(opts)
		if err != nil {
			return err
		}
		// no Kafka was selected, exit program
		if opts.name == "" {
			return nil
		}
	}

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

	var res *kafkamgmtclient.KafkaRequest
	ctx := context.Background()
	if opts.name != "" {
		res, _, err = kafka.GetKafkaByName(ctx, api.Kafka(), opts.name)
		if err != nil {
			return err
		}
	} else {
		res, _, err = kafka.GetKafkaByID(ctx, api.Kafka(), opts.id)
		if err != nil {
			return err
		}
	}

	// build Kafka config object from the response
	var kafkaConfig config.KafkaConfig = config.KafkaConfig{
		ClusterID: res.GetId(),
	}

	nameTmplEntry := localize.NewEntry("Name", res.GetName())
	cfg.Services.Kafka = &kafkaConfig
	if err := opts.Config.Save(cfg); err != nil {
		saveErrMsg := opts.localizer.LocalizeByID("kafka.use.error.saveError", nameTmplEntry)
		return fmt.Errorf("%v: %w", saveErrMsg, err)
	}

	logger.Info(opts.localizer.LocalizeByID("kafka.use.log.info.useSuccess", nameTmplEntry))

	return nil
}

func runInteractivePrompt(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger.Infoln(opts.localizer.LocalizeByID("common.log.debug.startingInteractivePrompt"))

	selectedKafka, err := kafka.InteractiveSelect(connection, logger)
	if err != nil {
		return err
	}

	opts.name = selectedKafka.GetName()

	return nil
}
