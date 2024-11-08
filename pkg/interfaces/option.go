package interfaces

import "github.com/liangweijiang/go-metric/pkg/config"

type Option interface {
	ApplyConfig(cfg *config.Config)
}
