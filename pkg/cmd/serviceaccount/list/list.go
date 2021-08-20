package list

import (
	"context"
	"encoding/json"

	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	"github.com/aerogear/charmil-host-example/pkg/cmdutil"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	output string
}

// svcAcctRow contains the properties used to
// populate the list of service accounts into a table row
type svcAcctRow struct {
	ID        string `json:"id" header:"ID"`
	Name      string `json:"name" header:"Name"`
	ClientID  string `json:"clientID" header:"Client ID"`
	Owner     string `json:"owner" header:"Owner"`
	CreatedAt string `json:"createdAt" header:"Created At"`
}

// NewListCommand creates a new command to list service accounts
func NewListCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("serviceAccount.list.cmd.use"),
		Short:   opts.localizer.LocalizeByID("serviceAccount.list.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("serviceAccount.list.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("serviceAccount.list.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.output != "" && !flagutil.IsValidInput(opts.output, flagutil.ValidOutputFormats...) {
				return flag.InvalidValueError("output", opts.output, flagutil.ValidOutputFormats...)
			}

			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", "", opts.localizer.LocalizeByID("serviceAccount.list.flag.output.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runList(opts *Options) (err error) {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	res, _, err := connection.API().ServiceAccount().GetServiceAccounts(context.Background()).Execute()
	if err != nil {
		return err
	}

	serviceaccounts := res.GetItems()
	if len(serviceaccounts) == 0 && opts.output == "" {
		logger.Info(opts.localizer.LocalizeByID("serviceAccount.list.log.info.noneFound"))
		return nil
	}

	outStream := opts.IO.Out
	switch opts.output {
	case dump.JSONFormat:
		data, _ := json.MarshalIndent(res, "", cmdutil.DefaultJSONIndent)
		_ = dump.JSON(outStream, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(res)
		_ = dump.YAML(outStream, data)
	default:
		rows := mapResponseItemsToRows(serviceaccounts)
		dump.Table(outStream, rows)
	}

	return nil
}

func mapResponseItemsToRows(svcAccts []kafkamgmtclient.ServiceAccountListItem) []svcAcctRow {
	rows := []svcAcctRow{}

	for _, sa := range svcAccts {
		row := svcAcctRow{
			ID:        sa.GetId(),
			Name:      sa.GetName(),
			ClientID:  sa.GetClientId(),
			Owner:     sa.GetOwner(),
			CreatedAt: sa.GetCreatedAt().String(),
		}

		rows = append(rows, row)
	}

	return rows
}
