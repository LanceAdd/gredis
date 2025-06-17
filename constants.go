package gredis

import (
	"time"
)

const (
	DefaultGroupName            = "default"
	frameCoreComponentNameRedis = "gf.core.component.redis"
	ConfigNodeNameRedis         = "redis"
)

const (
	defaultPoolMaxIdle     = 10
	defaultPoolMaxActive   = 100
	defaultPoolIdleTimeout = 10 * time.Second
	defaultPoolWaitTimeout = 10 * time.Second
	defaultPoolMaxLifeTime = 30 * time.Second
	defaultMaxRetries      = -1
)
