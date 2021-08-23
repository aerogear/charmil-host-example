package use

import (
	"context"
	"errors"

	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/serviceregistry"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil/core/utils/logging"
	srsmgmtv1 "github.com/redhat-developer/app-services-sdk-go/registrymgmt/apiv1/client"
	"github.com/spf13/cobra"
)

type Options struct {
	id          string
	name        string
	interactive bool

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewUseCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "use",
		Short:   f.Localizer.LocalizeByID("registry.cmd.use.shortDescription"),
		Long:    f.Localizer.LocalizeByID("registry.cmd.use.longDescription"),
		Example: f.Localizer.LocalizeByID("registry.cmd.use.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.id == "" && opts.name == "" {
				if !opts.IO.CanPrompt() {
					return errors.New(opts.localizer.LocalizeByID("registry.use.error.idOrNameRequired"))
				}
				opts.interactive = true
			}

			if opts.name != "" && opts.id != "" {
				return errors.New(opts.localizer.LocalizeByID("service.error.idAndNameCannotBeUsed"))
			}

			return runUse(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("registry.use.flag.id"))
	cmd.Flags().StringVar(&opts.name, "name", "", opts.localizer.LocalizeByID("registry.use.flag.name"))

	return cmd
}

func runUse(opts *Options) error {
	if opts.interactive {
		// run the use command interactively
		err := runInteractivePrompt(opts)
		if err != nil {
			return err
		}
		// no service was selected, exit program
		if opts.name == "" {
			return nil
		}
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	api := connection.API()

	var registry *srsmgmtv1.RegistryRest
	ctx := context.Background()
	if opts.name != "" {
		registry, _, err = serviceregistry.GetServiceRegistryByName(ctx, api.ServiceRegistryMgmt(), opts.name)
		if err != nil {
			return err
		}
	} else {
		registry, _, err = serviceregistry.GetServiceRegistryByID(ctx, api.ServiceRegistryMgmt(), opts.id)
		if err != nil {
			return err
		}
	}

	registryConfig := &config.ServiceRegistryConfig{
		InstanceID: registry.GetId(),
		Name:       *registry.Name,
	}

	nameTmplEntry := localize.NewEntry("Name", registry.GetName())
	opts.CfgHandler.Cfg.Services.ServiceRegistry = registryConfig

	logger.Info(opts.localizer.LocalizeByID("registry.use.log.info.useSuccess", nameTmplEntry))

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

	selectedRegistry, err := serviceregistry.InteractiveSelect(connection, logger)
	if err != nil {
		return err
	}

	opts.name = selectedRegistry.GetName()

	return nil
}
