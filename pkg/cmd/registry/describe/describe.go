package describe

import (
	"context"
	"errors"

	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/serviceregistry"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	srsmgmtv1 "github.com/redhat-developer/app-services-sdk-go/registrymgmt/apiv1/client"
	"github.com/spf13/cobra"
)

type Options struct {
	id           string
	name         string
	outputFormat string

	IO         *iostreams.IOStreams
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	localizer  localize.Localizer
}

// NewDescribeCommand describes a service instance, either by passing an `--id flag`
// or by using the service instance set in the config, if any
func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   f.Localizer.LocalizeByID("registry.cmd.describe.shortDescription"),
		Long:    f.Localizer.LocalizeByID("registry.cmd.describe.longDescription"),
		Example: f.Localizer.LocalizeByID("registry.cmd.describe.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			if len(args) > 0 {
				opts.name = args[0]
			}

			if opts.name != "" && opts.id != "" {
				return errors.New(opts.localizer.LocalizeByID("service.error.idAndNameCannotBeUsed"))
			}

			if opts.id != "" || opts.name != "" {
				return runDescribe(opts)
			}

			var registryConfig *config.ServiceRegistryConfig
			if opts.CfgHandler.Cfg.Services.ServiceRegistry == registryConfig || opts.CfgHandler.Cfg.Services.ServiceRegistry.InstanceID == "" {
				return errors.New(opts.localizer.LocalizeByID("registry.common.error.noServiceSelected"))
			}

			opts.id = opts.CfgHandler.Cfg.Services.ServiceRegistry.InstanceID

			return runDescribe(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("registry.cmd.flag.output.description"))
	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.LocalizeByID("registry.common.flag.id"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runDescribe(opts *Options) error {
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

	dump.PrintDataInFormat(opts.outputFormat, registry, opts.IO.Out)

	return nil
}
