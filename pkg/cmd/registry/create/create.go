package create

import (
	"context"
	"errors"

	"github.com/aerogear/charmil-host-example/pkg/serviceregistry"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/ams"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/connection"

	srsmgmtv1 "github.com/redhat-developer/app-services-sdk-go/registrymgmt/apiv1/client"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"

	"github.com/aerogear/charmil/core/utils/logging"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
)

type Options struct {
	name string

	outputFormat string
	autoUse      bool

	interactive bool

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

// NewCreateCommand creates a new command for creating registry.
func NewCreateCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   f.Localizer.LocalizeByID("registry.cmd.create.shortDescription"),
		Long:    f.Localizer.LocalizeByID("registry.cmd.create.longDescription"),
		Example: f.Localizer.LocalizeByID("registry.cmd.create.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.name = args[0]

				if err := serviceregistry.ValidateName(opts.name); err != nil {
					return err
				}
			}

			if !opts.IO.CanPrompt() && opts.name == "" {
				return errors.New(opts.localizer.LocalizeByID("registry.cmd.create.error.name.requiredWhenNonInteractive"))
			} else if opts.name == "" {
				opts.interactive = true
			}

			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			return runCreate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("registry.cmd.flag.output.description"))
	cmd.Flags().BoolVar(&opts.autoUse, "use", true, opts.localizer.LocalizeByID("registry.cmd.create.flag.use.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runCreate(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	var payload *srsmgmtv1.RegistryCreateRest
	if opts.interactive {
		logger.Infoln()

		payload, err = promptPayload(opts)
		if err != nil {
			return err
		}
	} else {
		payload = &srsmgmtv1.RegistryCreateRest{
			Name: &opts.name,
		}
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	// the user must have accepted the terms and conditions from the provider
	// before they can create a registry instance
	termsAccepted, termsURL, err := ams.CheckTermsAccepted(connection)
	if err != nil {
		return err
	}
	if !termsAccepted && termsURL != "" {
		logger.Info(opts.localizer.LocalizeByID("service.info.termsCheck", localize.NewEntry("TermsURL", termsURL)))
		return nil
	}

	logger.Info(opts.localizer.LocalizeByID("registry.cmd.create.info.action", localize.NewEntry("Name", payload.GetName())))

	response, _, err := connection.API().
		ServiceRegistryMgmt().
		CreateRegistry(context.Background()).
		RegistryCreateRest(*payload).
		Execute()
	if err != nil {
		return err
	}

	logger.Info(opts.localizer.LocalizeByID("registry.cmd.create.info.successMessage"))

	dump.PrintDataInFormat(opts.outputFormat, response, opts.IO.Out)

	registryConfig := &config.ServiceRegistryConfig{
		InstanceID: response.GetId(),
		Name:       response.GetName(),
	}

	if opts.autoUse {
		logger.Infoln("Auto-use is set, updating the current instance")
		opts.CfgHandler.Cfg.Services.ServiceRegistry = registryConfig
	} else {
		logger.Infoln("Auto-use is not set, skipping updating the current instance")
	}

	return nil
}

// Show a prompt to allow the user to interactively insert the data
func promptPayload(opts *Options) (payload *srsmgmtv1.RegistryCreateRest, err error) {
	if err != nil {
		return nil, err
	}

	// set type to store the answers from the prompt with defaults
	answers := struct {
		Name string
	}{}

	promptName := &survey.Input{
		Message: opts.localizer.LocalizeByID("registry.cmd.create.input.name.message"),
		Help:    opts.localizer.LocalizeByID("registry.cmd.create.input.name.help"),
	}

	err = survey.AskOne(promptName, &answers.Name, survey.WithValidator(serviceregistry.ValidateName))
	if err != nil {
		return nil, err
	}

	payload = &srsmgmtv1.RegistryCreateRest{
		Name: &answers.Name,
	}

	return payload, nil
}
