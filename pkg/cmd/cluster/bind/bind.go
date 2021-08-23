package bind

import (
	"context"
	"errors"

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

	forceCreationWithoutAsk bool
	ignoreContext           bool
	appName                 string
	selectedKafka           string

	forceOperator bool
	forceSDK      bool
	bindingName   string
}

func NewBindCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		CfgHandler: f.CfgHandler,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "bind",
		Short:   opts.localizer.LocalizeByID("cluster.bind.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("cluster.bind.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("cluster.bind.cmd.example"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.ignoreContext == true && !opts.IO.CanPrompt() {
				return errors.New(opts.localizer.LocalizeByID("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "ignore-context")))
			}
			if opts.appName == "" && !opts.IO.CanPrompt() {
				return errors.New(opts.localizer.LocalizeByID("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "appName")))
			}
			return runBind(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.kubeconfigLocation, "kubeconfig", "", "", opts.localizer.LocalizeByID("cluster.common.flag.kubeconfig.description"))
	cmd.Flags().StringVarP(&opts.appName, "app-name", "", "", opts.localizer.LocalizeByID("cluster.bind.flag.appName"))
	cmd.Flags().StringVarP(&opts.bindingName, "binding-name", "", "", opts.localizer.LocalizeByID("cluster.bind.flag.bindName"))
	cmd.Flags().BoolVarP(&opts.forceCreationWithoutAsk, "yes", "y", false, opts.localizer.LocalizeByID("cluster.common.flag.yes.description"))
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", opts.localizer.LocalizeByID("cluster.common.flag.namespace.description"))
	cmd.Flags().BoolVarP(&opts.ignoreContext, "ignore-context", "", false, opts.localizer.LocalizeByID("cluster.common.flag.ignoreContext.description"))
	cmd.Flags().BoolVarP(&opts.forceOperator, "force-operator", "", false, opts.localizer.LocalizeByID("cluster.bind.flag.forceOperator.description"))
	cmd.Flags().BoolVarP(&opts.forceSDK, "force-sdk", "", false, opts.localizer.LocalizeByID("cluster.bind.flag.forceSDK.description"))
	return cmd
}

func runBind(opts *Options) error {
	apiConnection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	// In future config will include Id's of other services
	if opts.CfgHandler.Cfg.Services.Kafka == nil || opts.ignoreContext {
		// nolint:govet
		selectedKafka, err := kafka.InteractiveSelect(apiConnection, logger)
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

	api := apiConnection.API()
	kafkaInstance, _, err := api.Kafka().GetKafkaById(context.Background(), opts.selectedKafka).Execute()
	if err != nil {
		return err
	}

	if kafkaInstance.Name == nil {
		return errors.New(opts.localizer.LocalizeByID("cluster.bind.error.emptyResponse"))
	}

	err = cluster.ExecuteServiceBinding(logger, opts.localizer, &cluster.ServiceBindingOptions{
		ServiceName:             kafkaInstance.GetName(),
		Namespace:               opts.namespace,
		AppName:                 opts.appName,
		ForceCreationWithoutAsk: opts.forceCreationWithoutAsk,
		ForceUseOperator:        opts.forceOperator,
		ForceUseSDK:             opts.forceSDK,
		BindingName:             opts.bindingName,
		BindAsFiles:             true,
	})

	return err
}
