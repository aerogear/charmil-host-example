// This file contains functions that help to manage visibility of early stage commands
package profile

import (
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
)

// Visual element displayed in help
func ApplyDevPreviewLabel(message string) string {
	return "[Preview] " + message
}

// Annotation used in templates and tools like documentation generation
func DevPreviewAnnotation() map[string]string {
	return map[string]string{"channel": "preview"}
}

// Check if preview is enabled
func DevPreviewEnabled(f *factory.Factory) bool {
	logger, err := f.Logger()
	if err != nil {
		logger.Info("Cannot determine status of dev preview. ", err)
		return false
	}

	return f.CfgHandler.Cfg.DevPreviewEnabled
}

// Enable dev preview
func EnableDevPreview(f *factory.Factory, enablement bool) (*config.Config, error) {
	logger, err := f.Logger()
	if err != nil {
		logger.Info(f.Localizer.LocalizeByID("profile.error.enablement"), err)
		return nil, err
	}

	f.CfgHandler.Cfg.DevPreviewEnabled = enablement

	if f.CfgHandler.Cfg.DevPreviewEnabled {
		logger.Info(f.Localizer.LocalizeByID("profile.status.devpreview.enabled"))
	} else {
		logger.Info(f.Localizer.LocalizeByID("profile.status.devpreview.disabled"))
	}
	return f.CfgHandler.Cfg, err
}
