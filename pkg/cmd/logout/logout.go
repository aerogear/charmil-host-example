// Package cluster contains commands for interacting with cluster logic of the service directly instead of through the
// REST API exposed via the serve command.
package logout

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

// NewLogoutCommand gets the command that's logs the current logged in user
func NewLogoutCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:   opts.localizer.LocalizeByID("logout.cmd.use"),
		Short: opts.localizer.LocalizeByID("logout.cmd.shortDescription"),
		Long:  opts.localizer.LocalizeByID("logout.cmd.longDescription"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runLogout(opts)
		},
	}
	return cmd
}

func runLogout(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	err = connection.Logout(context.TODO())

	if err != nil {
		return fmt.Errorf("%v: %w", opts.localizer.LocalizeByID("logout.error.unableToLogout"), err)
	}

	logger.Info(opts.localizer.LocalizeByID("logout.log.info.logoutSuccess"))

	return nil
}
