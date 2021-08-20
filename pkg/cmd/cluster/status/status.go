package status

import (
	"context"
	"fmt"

	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/cluster"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"

	"github.com/spf13/cobra"

	// Get all auth schemes
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/aerogear/charmil/core/utils/color"
	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	CfgHandler *config.CfgHandler
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	kubeconfig string
}

func NewStatusCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("cluster.status.cmd.use"),
		Short:   opts.localizer.LocalizeByID("cluster.status.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("cluster.status.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("cluster.status.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatus(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.kubeconfig, "kubeconfig", "", "", opts.localizer.LocalizeByID("cluster.common.flag.kubeconfig.description"))

	return cmd
}

func runStatus(opts *Options) error {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	clusterConn, err := cluster.NewKubernetesClusterConnection(connection, opts.CfgHandler, logger, opts.kubeconfig, opts.IO, opts.localizer)
	if err != nil {
		return err
	}

	var operatorStatus string
	// Add versioning in future
	isCRDInstalled, err := clusterConn.IsRhoasOperatorAvailableOnCluster(context.Background())
	if isCRDInstalled && err != nil {
		logger.Infoln(err)
	}

	if isCRDInstalled {
		operatorStatus = color.Success(opts.localizer.LocalizeByID("cluster.common.operatorInstalledMessage"))
	} else {
		operatorStatus = color.Error(opts.localizer.LocalizeByID("cluster.common.operatorNotInstalledMessage"))
	}

	currentNamespace, err := clusterConn.CurrentNamespace()
	if err != nil {
		return err
	}

	fmt.Fprintln(
		opts.IO.Out,
		opts.localizer.LocalizeByID("cluster.status.statusMessage",
			localize.NewEntry("Namespace", color.Info(currentNamespace)),
			localize.NewEntry("OperatorStatus", operatorStatus)),
	)

	return nil
}
