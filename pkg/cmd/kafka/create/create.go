package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	kafkamgmtclient "github.com/redhat-developer/app-services-sdk-go/kafkamgmt/apiv1/client"

	"github.com/aerogear/charmil/core/utils/localize"

	"github.com/aerogear/charmil-host-example/pkg/ams"
	"github.com/aerogear/charmil-host-example/pkg/cmd/flag"
	flagutil "github.com/aerogear/charmil-host-example/pkg/cmdutil/flags"
	"github.com/aerogear/charmil-host-example/pkg/connection"

	"github.com/aerogear/charmil-host-example/pkg/cloudprovider/cloudproviderutil"
	"github.com/aerogear/charmil-host-example/pkg/cloudregion/cloudregionutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aerogear/charmil-host-example/pkg/dump"
	pkgKafka "github.com/aerogear/charmil-host-example/pkg/kafka"
	"github.com/aerogear/charmil/core/utils/iostreams"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aerogear/charmil-host-example/internal/config"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/kafka/flags"
	"github.com/aerogear/charmil-host-example/pkg/cmdutil"

	"github.com/aerogear/charmil/core/utils/logging"
)

type Options struct {
	name     string
	provider string
	region   string
	multiAZ  bool

	outputFormat string
	autoUse      bool

	interactive bool

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

const (
	// default Kafka instance values
	defaultMultiAZ  = true
	defaultRegion   = "us-east-1"
	defaultProvider = "aws"
)

// NewCreateCommand creates a new command for creating kafkas.
func NewCreateCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,

		multiAZ: defaultMultiAZ,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.LocalizeByID("kafka.create.cmd.use"),
		Short:   opts.localizer.LocalizeByID("kafka.create.cmd.shortDescription"),
		Long:    opts.localizer.LocalizeByID("kafka.create.cmd.longDescription"),
		Example: opts.localizer.LocalizeByID("kafka.create.cmd.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				validator := &pkgKafka.Validator{
					Localizer:  opts.localizer,
					Connection: opts.Connection,
				}
				opts.name = args[0]

				if err := validator.ValidateName(opts.name); err != nil {
					return err
				}
			}

			if !opts.IO.CanPrompt() && opts.name == "" {
				return errors.New(opts.localizer.LocalizeByID("kafka.create.argument.name.error.requiredWhenNonInteractive"))
			} else if opts.name == "" {
				if opts.provider != "" || opts.region != "" {
					return errors.New(opts.localizer.LocalizeByID("kafka.create.argument.name.error.requiredWhenNonInteractive"))
				}
				opts.interactive = true
			}

			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			return runCreate(opts)
		},
	}

	cmd.Flags().StringVar(&opts.provider, flags.FlagProvider, "", opts.localizer.LocalizeByID("kafka.create.flag.cloudProvider.description"))
	cmd.Flags().StringVar(&opts.region, flags.FlagRegion, "", opts.localizer.LocalizeByID("kafka.create.flag.cloudRegion.description"))
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.LocalizeByID("kafka.common.flag.output.description"))
	cmd.Flags().BoolVar(&opts.autoUse, "use", true, opts.localizer.LocalizeByID("kafka.create.flag.autoUse.description"))

	_ = cmd.RegisterFlagCompletionFunc(flags.FlagProvider, func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return cmdutil.FetchCloudProviders(f)
	})

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

// nolint:funlen
func runCreate(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	cfg, err := opts.Config.Load()
	if err != nil {
		return err
	}

	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	// the user must have accepted the terms and conditions from the provider
	// before they can create a kafka instance
	termsAccepted, termsURL, err := ams.CheckTermsAccepted(connection)
	if err != nil {
		return err
	}
	if !termsAccepted && termsURL != "" {
		logger.Info(opts.localizer.LocalizeByID("service.info.termsCheck", localize.NewEntry("TermsURL", termsURL)))
		return nil
	}

	var payload *kafkamgmtclient.KafkaRequestPayload
	if opts.interactive {
		logger.Infoln()

		payload, err = promptKafkaPayload(opts)
		if err != nil {
			return err
		}

	} else {
		if opts.provider == "" {
			opts.provider = defaultProvider
		}
		if opts.region == "" {
			opts.region = defaultRegion
		}

		payload = &kafkamgmtclient.KafkaRequestPayload{
			Name:          opts.name,
			Region:        &opts.region,
			CloudProvider: &opts.provider,
			MultiAz:       &opts.multiAZ,
		}
	}

	logger.Info(opts.localizer.LocalizeByID("kafka.create.log.info.creatingKafka", localize.NewEntry("Name", payload.Name)))

	api := connection.API()

	a := api.Kafka().CreateKafka(context.Background())
	a = a.KafkaRequestPayload(*payload)
	a = a.Async(true)
	response, httpRes, err := a.Execute()

	if httpRes.StatusCode == 409 {
		return errors.New(opts.localizer.LocalizeByID("kafka.create.error.conflictError", localize.NewEntry("Name", payload.Name)))
	}

	if err != nil {
		return err
	}

	logger.Info(opts.localizer.LocalizeByID("kafka.create.info.successMessage", localize.NewEntry("Name", response.GetName())))

	switch opts.outputFormat {
	case dump.JSONFormat:
		data, _ := json.MarshalIndent(response, "", cmdutil.DefaultJSONIndent)
		_ = dump.JSON(opts.IO.Out, data)
	case dump.YAMLFormat, dump.YMLFormat:
		data, _ := yaml.Marshal(response)
		_ = dump.YAML(opts.IO.Out, data)
	}

	kafkaCfg := &config.KafkaConfig{
		ClusterID: response.GetId(),
	}

	if opts.autoUse {
		logger.Infoln("Auto-use is set, updating the current instance")
		cfg.Services.Kafka = kafkaCfg
		if err := opts.Config.Save(cfg); err != nil {
			return fmt.Errorf("%v: %w", opts.localizer.LocalizeByID("kafka.common.error.couldNotUseKafka"), err)
		}
	} else {
		logger.Infoln("Auto-use is not set, skipping updating the current instance")
	}

	return nil
}

// Show a prompt to allow the user to interactively insert the data for their Kafka
func promptKafkaPayload(opts *Options) (payload *kafkamgmtclient.KafkaRequestPayload, err error) {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return nil, err
	}

	api := connection.API()

	validator := &pkgKafka.Validator{
		Localizer:  opts.localizer,
		Connection: opts.Connection,
	}

	// set type to store the answers from the prompt with defaults
	answers := struct {
		Name          string
		Region        string
		MultiAZ       bool
		CloudProvider string
	}{
		MultiAZ: defaultMultiAZ,
	}

	promptName := &survey.Input{
		Message: opts.localizer.LocalizeByID("kafka.create.input.name.message"),
		Help:    opts.localizer.LocalizeByID("kafka.create.input.name.help"),
	}

	err = survey.AskOne(promptName, &answers.Name, survey.WithValidator(validator.ValidateName), survey.WithValidator(validator.ValidateNameIsAvailable))
	if err != nil {
		return nil, err
	}

	// fetch all cloud available providers
	cloudProviderResponse, _, err := api.Kafka().GetCloudProviders(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	cloudProviders := cloudProviderResponse.GetItems()
	cloudProviderNames := cloudproviderutil.GetEnabledNames(cloudProviders)

	cloudProviderPrompt := &survey.Select{
		Message: opts.localizer.LocalizeByID("kafka.create.input.cloudProvider.message"),
		Options: cloudProviderNames,
	}

	err = survey.AskOne(cloudProviderPrompt, &answers.CloudProvider)
	if err != nil {
		return nil, err
	}

	// get the selected provider type from the name selected
	selectedCloudProvider := cloudproviderutil.FindByName(cloudProviders, answers.CloudProvider)

	// nolint
	cloudRegionResponse, _, err := api.Kafka().GetCloudProviderRegions(context.Background(), selectedCloudProvider.GetId()).Execute()
	if err != nil {
		return nil, err
	}

	regions := cloudRegionResponse.GetItems()
	regionIDs := cloudregionutil.GetEnabledIDs(regions)

	regionPrompt := &survey.Select{
		Message: opts.localizer.LocalizeByID("kafka.create.input.cloudRegion.message"),
		Options: regionIDs,
		Help:    opts.localizer.LocalizeByID("kafka.create.input.cloudRegion.help"),
	}

	err = survey.AskOne(regionPrompt, &answers.Region)
	if err != nil {
		return nil, err
	}

	payload = &kafkamgmtclient.KafkaRequestPayload{
		Name:          answers.Name,
		Region:        &answers.Region,
		CloudProvider: &answers.CloudProvider,
		MultiAz:       &answers.MultiAZ,
	}

	return payload, nil
}
