package whoami

import (
	"fmt"

	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/auth/token"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/iostreams"
	"github.com/aerogear/charmil-host-example/pkg/localize"

	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	Config     config.IConfig
	Connection factory.ConnectionFunc
	IO         *iostreams.IOStreams
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewWhoAmICmd(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		IO:         f.IOStreams,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     f.Localizer.MustLocalize("whoami.cmd.use"),
		Short:   f.Localizer.MustLocalize("whoami.cmd.shortDescription"),
		Long:    f.Localizer.MustLocalize("whoami.cmd.longDescription"),
		Example: f.Localizer.MustLocalize("whoami.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCmd(opts)
		},
	}

	return cmd
}

func runCmd(opts *Options) (err error) {
	cfg, err := opts.Config.Load()
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	_, err = opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	accessTkn, _ := token.Parse(cfg.AccessToken)

	tknClaims, _ := token.MapClaims(accessTkn)

	userName, ok := tknClaims["preferred_username"]

	if ok {
		fmt.Fprintln(opts.IO.Out, userName)
	} else {
		logger.Info(opts.localizer.MustLocalize("whoami.log.info.tokenHasNoUsername"))
	}

	return nil
}
