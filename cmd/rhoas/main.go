package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/doc"
	"github.com/aerogear/charmil-host-example/pkg/localesettings"
	"github.com/aerogear/charmil/core/utils/localize"
	"golang.org/x/text/language"

	"github.com/aerogear/charmil-host-example/internal/build"

	"github.com/aerogear/charmil-host-example/pkg/cmd/debug"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/cmd/root"
	"github.com/spf13/cobra"
)

var generateDocs = os.Getenv("GENERATE_DOCS") == "true"

func main() {

	cfg := &config.Config{
		Services: &config.ServiceConfigMap{
			Kafka:           &config.KafkaConfig{},
			ServiceRegistry: &config.ServiceRegistryConfig{},
		},
	}

	cfgHandler, err := config.NewHandler(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	locConfig := &localize.Config{
		Language: &language.English,
		Files:    localesettings.DefaultLocales,
		Format:   "toml",
	}

	localizer, err := localize.New(locConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	buildVersion := build.Version
	cmdFactory := factory.New(buildVersion, localizer, cfgHandler)
	logger, err := cmdFactory.Logger()
	if err != nil {
		fmt.Println(cmdFactory.IOStreams.ErrOut, err)
		os.Exit(1)
	}

	if err = cmdFactory.CfgHandler.Load(); err != nil {
		fmt.Println(cmdFactory.IOStreams.ErrOut, err)
		os.Exit(1)
	}

	rootCmd := root.NewRootCommand(cmdFactory, buildVersion)

	rootCmd.InitDefaultHelpCmd()

	if generateDocs {
		generateDocumentation(rootCmd)
		os.Exit(0)
	}

	err = rootCmd.Execute()
	if err != nil {
		logger.Error(wrapErrorf(err, localizer))
		build.CheckForUpdate(context.Background(), logger, localizer)
		os.Exit(1)
	}

	if debug.Enabled() {
		build.CheckForUpdate(context.Background(), logger, localizer)
	}

	if err = cmdFactory.CfgHandler.Save(); err != nil {
		fmt.Println(cmdFactory.IOStreams.ErrOut, err)
		os.Exit(1)
	}
}

/**
* Generates documentation files
 */
func generateDocumentation(rootCommand *cobra.Command) {
	fmt.Fprint(os.Stderr, "Generating docs.\n\n")
	filePrepender := func(filename string) string {
		return ""
	}

	rootCommand.DisableAutoGenTag = true

	linkHandler := func(s string) string { return s }

	err := doc.GenAsciidocTreeCustom(rootCommand, "./docs/commands", filePrepender, linkHandler)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func wrapErrorf(err error, localizer localize.Localizer) error {
	return fmt.Errorf("Error: %w. %v", err, localizer.LocalizeByID("common.log.error.verboseModeHint"))
}
