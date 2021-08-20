package config

import (
	"errors"
	"strconv"

	"github.com/aerogear/charmil-host-example/pkg/profile"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/logging"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/spf13/cobra"
)

type Options struct {
	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewConfigCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		CfgHandler: f.CfgHandler,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "config",
		Short:   opts.localizer.LocalizeByID("config.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("config.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("config.cmd.example"),
	}

	devPreview := &cobra.Command{
		Use:       "dev-preview",
		Short:     opts.localizer.LocalizeByID("devpreview.cmd.shortDescription"),
		Long:      opts.localizer.LocalizeByID("devpreview.cmd.longDescription"),
		Example:   opts.localizer.LocalizeByID("devpreview.cmd.example"),
		ValidArgs: []string{"true", "false"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			devPreview, err := strconv.ParseBool(args[0])
			if err != nil {
				return errors.New(opts.localizer.LocalizeByID("devpreview.error.enablement"))
			}
			_, err = profile.EnableDevPreview(f, devPreview)
			return err
		},
	}
	cmd.AddCommand(devPreview)
	return cmd
}
