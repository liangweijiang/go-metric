package interfaces

import "go-mertric/pkg/config"

type Option interface {
	ApplyConfig(cfg *config.Config)
}
