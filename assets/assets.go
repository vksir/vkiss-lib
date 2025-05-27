package assets

import (
	_ "embed"
)

//go:embed config.toml
var DefaultConfig string
