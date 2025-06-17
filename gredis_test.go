package gredis

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var (
	ctx    = gctx.GetInitCtx()
	config = &Config{
		Address: `127.0.0.1:6379`,
		Db:      1,
	}
	client = New(config)
)

func TestRedis(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		nx := client.SetEx(ctx, "flag", "OK", 10*time.Minute)
		glog.Info(ctx, nx.Val())
	})
}

func TestRedisPipeline(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		var incr *redis.IntCmd
		client.Set(ctx, "number", 0, 10*time.Minute)
		cmds, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for i := 0; i < 10; i++ {
				incr = pipe.Incr(ctx, "number")
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		t.AssertEQ(len(cmds), 10)
		t.AssertEQ(incr.Val(), int64(10))
	})
}
