package main

import (
	"context"
	"fmt"
	"os"

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
	cmdFactory := factory.New(build.Version, localizer)
	logger, err := cmdFactory.Logger()
	if err != nil {
		fmt.Println(cmdFactory.IOStreams.ErrOut, err)
		os.Exit(1)
	}

	cfgFile, err := cmdFactory.Config.Load()
	if err != nil {
		fmt.Println(cmdFactory.IOStreams.ErrOut, err)
		os.Exit(1)
	}

	fmt.Println(cfgFile)

	rootCmd := root.NewRootCommand(cmdFactory, buildVersion)

	rootCmd.InitDefaultHelpCmd()

	if generateDocs {
		generateDocumentation(rootCmd)
		os.Exit(0)
	}

	if err = cmdFactory.Config.Save(cfgFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = rootCmd.Execute()
	if err == nil {
		if debug.Enabled() {
			build.CheckForUpdate(context.Background(), logger, localizer)
		}
		return
	}

	if err != nil {
		logger.Error(wrapErrorf(err, localizer))
		build.CheckForUpdate(context.Background(), logger, localizer)
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
