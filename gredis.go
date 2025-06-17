package gredis

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/redis/go-redis/v9"
)

var (
	localInstances = gmap.NewStrAnyMap(true)
)

func Redis(name ...string) redis.UniversalClient {
	var (
		ctx   = context.Background()
		group = DefaultGroupName
	)
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameRedis, group)
	result := localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		if c, ok := GetConfig(group); ok {
			return New(c)
		}
		config, err := loadConfigFromGlobal(ctx, group)
		if err != nil {
			glog.Errorf(ctx, "failed to load redis config for group: %q: %+v", group, err)
		}
		if config == nil {
			glog.Printf(context.Background(), "missing configuration for redis group: %q", group)
			return nil
		}
		localConfigMap.GetOrSetFuncLock(group, func() interface{} {
			return config
		})
		return New(config)
	})
	if result == nil {
		return nil
	}
	return result.(redis.UniversalClient)
}

func New(config *Config) redis.UniversalClient {
	fillWithDefaultConfiguration(config)
	opts := &redis.UniversalOptions{
		Addrs:            gstr.SplitAndTrim(config.Address, ","),
		Username:         config.User,
		Password:         config.Pass,
		SentinelUsername: config.SentinelUser,
		SentinelPassword: config.SentinelPass,
		DB:               config.Db,
		MaxRetries:       defaultMaxRetries,
		PoolSize:         config.MaxActive,
		MinIdleConns:     config.MinIdle,
		MaxIdleConns:     config.MaxIdle,
		ConnMaxLifetime:  config.MaxConnLifetime,
		ConnMaxIdleTime:  config.IdleTimeout,
		PoolTimeout:      config.WaitTimeout,
		DialTimeout:      config.DialTimeout,
		ReadTimeout:      config.ReadTimeout,
		WriteTimeout:     config.WriteTimeout,
		MasterName:       config.MasterName,
		TLSConfig:        config.TLSConfig,
		Protocol:         config.Protocol,
	}
	var client redis.UniversalClient
	if opts.MasterName != "" {
		redisSentinel := opts.Failover()
		redisSentinel.ReplicaOnly = config.SlaveOnly
		client = redis.NewFailoverClient(redisSentinel)
	} else if len(opts.Addrs) > 1 || config.Cluster {
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		client = redis.NewClient(opts.Simple())
	}
	return client
}

func loadConfigFromGlobal(ctx context.Context, group string) (*Config, error) {
	if !gcfg.Instance().Available(ctx) {
		return nil, gerror.New("global config is not available")
	}
	configMap, err := gcfg.Instance().Data(ctx)
	if err != nil {
		return nil, err
	}
	if _, v := gutil.MapPossibleItemByKey(configMap, ConfigNodeNameRedis); v != nil {
		configMap = gconv.Map(v)
	}
	if len(configMap) == 0 {
		return nil, nil
	}
	if v, ok := configMap[group]; ok {
		return ConfigFromMap(gconv.Map(v))
	}
	return nil, nil
}

func Clear() {
	localConfigMap.LockFunc(func(m map[string]interface{}) {
		m = make(map[string]interface{})
		localInstances.Clear()
	})
	glog.Info(context.Background(), "redis config cleared")
}
