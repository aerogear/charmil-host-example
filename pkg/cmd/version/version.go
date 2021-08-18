package version

import (
	"context"
	"fmt"

	"github.com/aerogear/charmil-host-example/internal/build"
	"github.com/aerogear/charmil-host-example/pkg/cmd/debug"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/iostreams"
	"github.com/aerogear/charmil-host-example/pkg/localize"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	IO        *iostreams.IOStreams
	Logger    func() (logging.Logger, error)
	localizer localize.Localizer
}

func NewVersionCmd(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:        f.IOStreams,
		Logger:    f.Logger,
		localizer: f.Localizer,
	}

	cmd := &cobra.Command{
		Use:    opts.localizer.MustLocalize("version.cmd.use"),
		Short:  opts.localizer.MustLocalize("version.cmd.shortDescription"),
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCmd(opts)
		},
	}

	return cmd
}

func runCmd(opts *Options) (err error) {
	fmt.Fprintln(opts.IO.Out, opts.localizer.MustLocalize("version.cmd.outputText", localize.NewEntry("Version", build.Version)))

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	// debug mode checks this for a version update also.
	// so we check if is enabled first so as not to print it twice
	if !debug.Enabled() {
		build.CheckForUpdate(context.Background(), logger, opts.localizer)
	}
	return nil
}
