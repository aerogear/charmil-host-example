package create

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aerogear/charmil-host-example/pkg/serviceaccount/validation"
	"github.com/aerogear/charmil/core/utils/localize"
	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"

	"github.com/aerogear/charmil-host-example/pkg/connection"

	"github.com/AlecAivazis/survey/v2"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/serviceaccount/credentials"
	"github.com/aerogear/charmil/core/utils/iostreams"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer

	fileFormat  string
	overwrite   bool
	name        string
	description string
	filename    string

	interactive bool
}

// NewCreateCommand creates a new command to create service accounts
func NewCreateCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("serviceAccount.create.cmd.use"),
		Short:   opts.localizer.LocalizeByID("serviceAccount.create.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("serviceAccount.create.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("serviceAccount.create.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			if !opts.IO.CanPrompt() && opts.name == "" {
				return errors.New(opts.localizer.LocalizeByID("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "name")))
			} else if opts.name == "" && opts.description == "" {
				opts.interactive = true
			}

			if !opts.interactive {

				validator := &validation.Validator{
					Localizer: opts.localizer,
				}

				if opts.fileFormat == "" {
					return errors.New(opts.localizer.LocalizeByID("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "file-format")))
				}

				if err = validator.ValidateName(opts.name); err != nil {
					return err
				}
				if err = validator.ValidateDescription(opts.description); err != nil {
					return err
				}
			}

			// check that a valid --file-format flag value is used
			validOutput := flagutil.IsValidInput(opts.fileFormat, flagutil.CredentialsOutputFormats...)
			if !validOutput && opts.fileFormat != "" {
				return flag.InvalidValueError("file-format", opts.fileFormat, flagutil.CredentialsOutputFormats...)
			}

			return runCreate(opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", opts.localizer.LocalizeByID("serviceAccount.create.flag.name.description"))
	cmd.Flags().StringVar(&opts.description, "description", "", opts.localizer.LocalizeByID("serviceAccount.create.flag.description.description"))
	cmd.Flags().BoolVar(&opts.overwrite, "overwrite", false, opts.localizer.LocalizeByID("serviceAccount.common.flag.overwrite.description"))
	cmd.Flags().StringVar(&opts.filename, "file-location", "", opts.localizer.LocalizeByID("serviceAccount.common.flag.fileLocation.description"))
	cmd.Flags().StringVar(&opts.fileFormat, "file-format", "", opts.localizer.LocalizeByID("serviceAccount.common.flag.fileFormat.description"))

	flagutil.EnableStaticFlagCompletion(cmd, "file-format", flagutil.CredentialsOutputFormats)

	return cmd
}

// nolint:funlen
func runCreate(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	if opts.interactive {
		// run the create command interactively
		err = runInteractivePrompt(opts)
		if err != nil {
			return err
		}
	} else if opts.filename == "" {
		// obtain the absolute path to where credentials will be saved
		opts.filename = credentials.GetDefaultPath(opts.fileFormat)
	}

	// If the credentials file already exists, and the --overwrite flag is not set then return an error
	// indicating that the user should explicitly request overwriting of the file
	_, err = os.Stat(opts.filename)
	if err == nil && !opts.overwrite {
		return errors.New(opts.localizer.LocalizeByID("serviceAccount.common.error.credentialsFileAlreadyExists", localize.NewEntry("FilePath", opts.filename)))
	}

	// create the service account
	serviceAccountPayload := &kafkamgmtclient.ServiceAccountRequest{Name: opts.name, Description: &opts.description}

	a := connection.API().ServiceAccount().CreateServiceAccount(context.Background())
	a = a.ServiceAccountRequest(*serviceAccountPayload)
	serviceacct, _, err := a.Execute()
	if err != nil {
		return err
	}

	logger.Info(opts.localizer.LocalizeByID("serviceAccount.create.log.info.createdSuccessfully", localize.NewEntry("ID", serviceacct.GetId()), localize.NewEntry("Name", serviceacct.GetName())))

	creds := &credentials.Credentials{
		ClientID:     serviceacct.GetClientId(),
		ClientSecret: serviceacct.GetClientSecret(),
	}

	// save the credentials to a file
	err = credentials.Write(opts.fileFormat, opts.filename, creds)
	if err != nil {
		return fmt.Errorf("%v: %w", opts.localizer.LocalizeByID("serviceAccount.common.error.couldNotSaveCredentialsFile"), err)
	}

	logger.Info(opts.localizer.LocalizeByID("serviceAccount.common.log.info.credentialsSaved", localize.NewEntry("FilePath", opts.filename)))

	return nil
}

func runInteractivePrompt(opts *Options) (err error) {
	_, err = opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	validator := &validation.Validator{
		Localizer: opts.localizer,
	}

	logger.Infoln(opts.localizer.LocalizeByID("common.log.debug.startingInteractivePrompt"))

	promptName := &survey.Input{
		Message: opts.localizer.LocalizeByID("serviceAccount.create.input.name.message"),
		Help:    opts.localizer.LocalizeByID("serviceAccount.create.input.name.help"),
	}

	err = survey.AskOne(promptName, &opts.name, survey.WithValidator(survey.Required), survey.WithValidator(validator.ValidateName))
	if err != nil {
		return err
	}

	// if the --file-format flag was not used, ask in the prompt
	if opts.fileFormat == "" {
		logger.Infoln(opts.localizer.LocalizeByID("serviceAccount.common.log.debug.interactive.fileFormatNotSet"))

		fileFormatPrompt := &survey.Select{
			Message: opts.localizer.LocalizeByID("serviceAccount.create.input.fileFormat.message"),
			Help:    opts.localizer.LocalizeByID("serviceAccount.create.input.fileFormat.help"),
			Options: flagutil.CredentialsOutputFormats,
			Default: "env",
		}

		err = survey.AskOne(fileFormatPrompt, &opts.fileFormat)
		if err != nil {
			return err
		}
	}

	opts.filename, opts.overwrite, err = credentials.ChooseFileLocation(opts.fileFormat, opts.filename, opts.overwrite)
	if err != nil {
		return err
	}

	promptDescription := &survey.Multiline{
		Message: opts.localizer.LocalizeByID("serviceAccount.create.input.description.message"),
		Help:    opts.localizer.LocalizeByID("serviceAccount.create.flag.description.description"),
	}

	err = survey.AskOne(promptDescription, &opts.description, survey.WithValidator(validator.ValidateDescription))
	if err != nil {
		return err
	}

	logger.Info(opts.localizer.LocalizeByID("serviceAccount.create.log.info.creating", localize.NewEntry("Name", opts.name)))

	return nil
}
