package connect

import (
	"context"
	"errors"

	"github.com/aerogear/charmil-host-example/internal/build"
	"github.com/aerogear/charmil-host-example/pkg/cluster"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil-host-example/pkg/kafka"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/spf13/cobra"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	CfgHandler *config.CfgHandler
	Connection func(connectionCfg *connection.Config) (connection.Connection, error)
	Logger     func() (logging.Logger, error)
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	kubeconfigLocation string
	namespace          string

	offlineAccessToken      string
	forceCreationWithoutAsk bool
	ignoreContext           bool
	selectedKafka           string
}

func NewConnectCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("cluster.connect.cmd.use"),
		Short:   opts.localizer.LocalizeByID("cluster.connect.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("cluster.connect.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("cluster.connect.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.ignoreContext == true && !opts.IO.CanPrompt() {
				return errors.New(opts.localizer.LocalizeByID("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "ignore-context")))
			}
			return runConnect(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.kubeconfigLocation, "kubeconfig", "", "", opts.localizer.LocalizeByID("cluster.common.flag.kubeconfig.description"))
	cmd.Flags().StringVarP(&opts.offlineAccessToken, "token", "", "", opts.localizer.LocalizeByID("cluster.common.flag.offline.token.description", localize.NewEntry("OfflineTokenURL", build.OfflineTokenURL)))
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", opts.localizer.LocalizeByID("cluster.common.flag.namespace.description"))
	cmd.Flags().BoolVarP(&opts.forceCreationWithoutAsk, "yes", "y", false, opts.localizer.LocalizeByID("cluster.common.flag.yes.description"))
	cmd.Flags().BoolVarP(&opts.ignoreContext, "ignore-context", "", false, opts.localizer.LocalizeByID("cluster.common.flag.ignoreContext.description"))

	return cmd
}

func runConnect(opts *Options) error {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	clusterConn, err := cluster.NewKubernetesClusterConnection(connection, opts.CfgHandler, logger, opts.kubeconfigLocation, opts.IO, opts.localizer)
	if err != nil {
		return err
	}

	// In future config will include Id's of other services
	if opts.CfgHandler.Cfg.Services.Kafka == nil || opts.ignoreContext {
		// nolint
		selectedKafka, err := kafka.InteractiveSelect(connection, logger)
		if err != nil {
			return err
		}
		if selectedKafka == nil {
			return nil
		}
		opts.selectedKafka = selectedKafka.GetId()
	} else {
		opts.selectedKafka = opts.CfgHandler.Cfg.Services.Kafka.ClusterID
	}

	arguments := &cluster.ConnectArguments{
		OfflineAccessToken:      opts.offlineAccessToken,
		ForceCreationWithoutAsk: opts.forceCreationWithoutAsk,
		IgnoreContext:           opts.ignoreContext,
		SelectedKafka:           opts.selectedKafka,
		Namespace:               opts.namespace,
	}

	err = clusterConn.Connect(context.Background(), arguments)
	if err != nil {
		return err
	}

	return nil
}
