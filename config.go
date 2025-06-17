package gredis

import (
	"crypto/tls"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

var (
	// Configuration groups.
	localConfigMap = gmap.NewStrAnyMap(true)
)

type Config struct {
	// Address It supports single and cluster redis server. Multiple addresses joined with char ','. Eg: 192.168.1.1:6379, 192.168.1.2:6379.
	Address         string        `json:"address"`
	Db              int           `json:"db"`              // Redis db.
	User            string        `json:"user"`            // Username for AUTH.
	Pass            string        `json:"pass"`            // Password for AUTH.
	SentinelUser    string        `json:"sentinel_user"`   // Username for sentinel AUTH.
	SentinelPass    string        `json:"sentinel_pass"`   // Password for sentinel AUTH.
	MinIdle         int           `json:"minIdle"`         // Minimum number of connections allowed to be idle (default is 0)
	MaxIdle         int           `json:"maxIdle"`         // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int           `json:"maxActive"`       // Maximum number of connections limit (default is 0 means no limit).
	MaxConnLifetime time.Duration `json:"maxConnLifetime"` // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	IdleTimeout     time.Duration `json:"idleTimeout"`     // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	WaitTimeout     time.Duration `json:"waitTimeout"`     // Timed out duration waiting to get a connection from the connection pool.
	DialTimeout     time.Duration `json:"dialTimeout"`     // Dial connection timeout for TCP.
	ReadTimeout     time.Duration `json:"readTimeout"`     // Read timeout for TCP. DO NOT set it if not necessary.
	WriteTimeout    time.Duration `json:"writeTimeout"`    // Write timeout for TCP.
	MasterName      string        `json:"masterName"`      // Used in Redis Sentinel mode.
	TLS             bool          `json:"tls"`             // Specifies whether TLS should be used when connecting to the server.
	TLSSkipVerify   bool          `json:"tlsSkipVerify"`   // Disables server name verification when connecting over TLS.
	TLSConfig       *tls.Config   `json:"-"`               // TLS Config to use. When set TLS will be negotiated.
	SlaveOnly       bool          `json:"slaveOnly"`       // Route all commands to slave read-only nodes.
	Cluster         bool          `json:"cluster"`         // Specifies whether cluster mode be used.
	Protocol        int           `json:"protocol"`        // Specifies the RESP version (Protocol 2 or 3.)
}

func GetConfig(name ...string) (config *Config, ok bool) {
	group := DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	if v := localConfigMap.Get(group); v != nil {
		return v.(*Config), true
	}
	return &Config{}, false
}

func ConfigFromMap(m map[string]interface{}) (config *Config, err error) {
	config = &Config{}
	if err = gconv.Scan(m, config); err != nil {
		err = gerror.NewCodef(gcode.CodeInvalidConfiguration, `invalid redis configuration: %#v`, m)
	}
	if config.DialTimeout < time.Second {
		config.DialTimeout = config.DialTimeout * time.Second
	}
	if config.WaitTimeout < time.Second {
		config.WaitTimeout = config.WaitTimeout * time.Second
	}
	if config.WriteTimeout < time.Second {
		config.WriteTimeout = config.WriteTimeout * time.Second
	}
	if config.ReadTimeout < time.Second {
		config.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.IdleTimeout < time.Second {
		config.IdleTimeout = config.IdleTimeout * time.Second
	}
	if config.MaxConnLifetime < time.Second {
		config.MaxConnLifetime = config.MaxConnLifetime * time.Second
	}
	if config.Protocol != 2 && config.Protocol != 3 {
		config.Protocol = 3
	}
	return
}

func fillWithDefaultConfiguration(config *Config) {
	if config.MaxIdle == 0 {
		config.MaxIdle = defaultPoolMaxIdle
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = defaultPoolMaxActive
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultPoolIdleTimeout
	}
	if config.WaitTimeout == 0 {
		config.WaitTimeout = defaultPoolWaitTimeout
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = defaultPoolMaxLifeTime
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = -1
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = -1
	}
	if config.TLSConfig == nil && config.TLS {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.TLSSkipVerify,
		}
	}
}
