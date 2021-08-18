package describe

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Options struct {
	id           string
	outputFormat string

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	localizer  localize.Localizer
}

func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("serviceAccount.describe.cmd.use"),
		Short:   opts.localizer.LocalizeByID("serviceAccount.describe.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("serviceAccount.describe.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("serviceAccount.describe.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			return runDescribe(opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("serviceAccount.describe.flag.id.description"))
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("serviceAccount.common.flag.output.description"))

	_ = cmd.MarkFlagRequired("id")

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runDescribe(opts *Options) error {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	api := connection.API()

	res, httpRes, err := api.ServiceAccount().GetServiceAccountById(context.Background(), opts.id).Execute()
	if err != nil {
		if httpRes == nil {
			return err
		}

		switch httpRes.StatusCode {
		case 404:
			return errors.New(opts.localizer.LocalizeByID("serviceAccount.common.error.notFoundError", localize.NewEntry("ID", opts.id)))
		default:
			return err
		}
	}

	switch opts.outputFormat {
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(res)
		_ = dump.YAML(opts.IO.Out, data)
	default:
		data, _ := json.Marshal(res)
		_ = dump.JSON(opts.IO.Out, data)
	}

	return nil
}
