package profile

import (
	"testing"

	"github.com/aerogear/charmil-host-example/internal/mockutil"
	"github.com/aerogear/charmil-host-example/pkg/cmd/factory"
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/localesettings"
	"github.com/aerogear/charmil/core/utils/localize"
	"golang.org/x/text/language"
)

func TestEnableDevPreviewConfig(t *testing.T) {
	locConfig := &localize.Config{
		Language: &language.English,
		Files:    localesettings.DefaultLocales,
		Format:   "toml",
	}

	localizer, _ := localize.New(locConfig)
	testVal := true
	factoryObj := factory.New("dev", localizer)

	factoryObj.Config = mockutil.NewConfigMock(&config.Config{})
	config, err := EnableDevPreview(factoryObj, testVal)
	if config.DevPreviewEnabled == false {
		t.Errorf("TestEnableDevPreviewConfig config = %v, want %v", config.DevPreviewEnabled, true)
	}
	if err != nil {
		t.Errorf("TestEnableDevPreviewConfig error = %v, want %v", err, nil)
	}

	testVal = false
	config, err = EnableDevPreview(factoryObj, testVal)
	if config.DevPreviewEnabled == true {
		t.Errorf("TestEnableDevPreviewConfig() config.DevPreviewEnabled = %v, want %v", config.DevPreviewEnabled, false)
	}
	if err != nil {
		t.Errorf("TestEnableDevPreviewConfig error = %v, want %v", err, nil)
	}

}

func TestDevPreviewEnabled(t *testing.T) {
	locConfig := &localize.Config{
		Language: &language.English,
		Files:    localesettings.DefaultLocales,
		Format:   "toml",
	}

	localizer, _ := localize.New(locConfig)
	factoryObj := factory.New("dev", localizer)
	factoryObj.Config = mockutil.NewConfigMock(&config.Config{})
	testVal := false
	_, err := EnableDevPreview(factoryObj, testVal)
	if err != nil {
		t.Errorf("TestEnableDevPreviewConfig error = %v, want %v", err, nil)
	}
	enabled := DevPreviewEnabled(factoryObj)

	if enabled {
		t.Errorf("TestEnableDevPreviewConfig enabled = %v, want %v", enabled, false)
	}

	testVal = true
	_, _ = EnableDevPreview(factoryObj, testVal)

	enabled = DevPreviewEnabled(factoryObj)
	if !enabled {
		t.Errorf("TestEnableDevPreviewConfig enabled = %v, want %v", enabled, true)
	}
}
