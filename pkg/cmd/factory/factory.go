package factory

import (
	"github.com/aerogear/charmil-host-example/pkg/config"
	"github.com/aerogear/charmil-host-example/pkg/connection"
	"github.com/aerogear/charmil/core/utils/iostreams"
	"github.com/aerogear/charmil/core/utils/localize"
	"github.com/aerogear/charmil/core/utils/logging"
)

// Factory is an abstract type which provides access to
// the root configuration and connections for the CLI
type Factory struct {
	// Type which defines the streams for the CLI
	IOStreams *iostreams.IOStreams

	// Creates a connection to the API
	Connection ConnectionFunc

	// Returns a logger to create leveled logs in the application
	Logger func() (logging.Logger, error)

	// Localizer provides text to the commands
	Localizer localize.Localizer

	// CfgHandler provides the fields required for managing config
	CfgHandler *config.CfgHandler
}

type ConnectionFunc func(cfg *connection.Config) (connection.Connection, error)
