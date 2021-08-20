package delete

import (
	"context"
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/serviceaccount/validation"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer

	id    string
	force bool
}

// NewDeleteCommand creates a new command to delete a service account
func NewDeleteCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("serviceAccount.delete.cmd.use"),
		Short:   opts.localizer.LocalizeByID("serviceAccount.delete.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("serviceAccount.delete.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("serviceAccount.delete.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !opts.IO.CanPrompt() && !opts.force {
				return flag.RequiredWhenNonInteractiveError("yes")
			}

			validator := &validation.Validator{
				Localizer: opts.localizer,
			}

			validID := validator.ValidateUUID(opts.id)
			if validID != nil {
				return validID
			}

			return runDelete(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("serviceAccount.delete.flag.id.description"))
	cmd.Flags().BoolVarP(&opts.force, "yes", "y", false, opts.localizer.LocalizeByID("serviceAccount.delete.flag.yes.description"))

	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func runDelete(opts *Options) (err error) {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	_, httpRes, err := connection.API().ServiceAccount().GetServiceAccountById(context.Background(), opts.id).Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		if httpRes.StatusCode == 404 {
			return errors.New(opts.localizer.LocalizeByID("serviceAccount.common.error.notFoundError", localize.NewEntry("ID", opts.id)))
		}
	}

	if !opts.force {
		var confirmDelete bool
		promptConfirmDelete := &survey.Confirm{
			Message: opts.localizer.LocalizeByID("serviceAccount.delete.input.confirmDelete.message", localize.NewEntry("ID", opts.id)),
		}

		err = survey.AskOne(promptConfirmDelete, &confirmDelete)
		if err != nil {
			return err
		}

		if !confirmDelete {
			logger.Infoln(opts.localizer.LocalizeByID("serviceAccount.delete.log.debug.deleteNotConfirmed"))
			return nil
		}
	}

	return deleteServiceAccount(opts)
}

func deleteServiceAccount(opts *Options) error {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	_, httpRes, err := connection.API().ServiceAccount().DeleteServiceAccountById(context.Background(), opts.id).Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		switch httpRes.StatusCode {
		case 403:
			return errors.New(opts.localizer.LocalizeByID("serviceAccount.common.error.forbidden", localize.NewEntry("Operation", "delete")))
		case 500:
			return errors.New(opts.localizer.LocalizeByID("serviceAccount.common.error.internalServerError"))
		default:
			return err
		}
	}

	logger.Info(opts.localizer.LocalizeByID("serviceAccount.delete.log.info.deleteSuccess"))

	return nil
}
