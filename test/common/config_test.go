package testing

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	CommonCfg "github.com/webalytic.go/common/config"
)

func TestLoadRedisConfigFromDefaultConfig(t *testing.T) {
	v := CommonCfg.Init()
	redisConfig := CommonCfg.RedisConfig{V: v}
	assert.Equal(t, redisConfig.Port(), 6379)
	assert.Equal(t, redisConfig.Host(), "0.0.0.0")
}

func TestLoadRedisEnvShouldBePriorForCollector(t *testing.T) {
	v := CommonCfg.Init()
	redisConfig := CommonCfg.RedisConfig{V: v}
	v.SetEnvPrefix("COLLECTOR")
	os.Setenv("COLLECTOR_REDIS_PORT", "1234")
	os.Setenv("COLLECTOR_REDIS_HOST", "127.0.0.1")
	assert.Equal(t, 1234, redisConfig.Port())
	assert.Equal(t, "127.0.0.1", redisConfig.Host())
}

func TestLoadRedisEnvShouldBePriorForHandler(t *testing.T) {
	v := CommonCfg.Init()
	redisConfig := CommonCfg.RedisConfig{V: v}
	v.SetEnvPrefix("HANDLER")
	os.Setenv("HANDLER_REDIS_PORT", "1234")
	os.Setenv("HANDLER_REDIS_HOST", "127.0.0.1")
	assert.Equal(t, 1234, redisConfig.Port())
	assert.Equal(t, "127.0.0.1", redisConfig.Host())
}

func TestNoDefaultFileFound(t *testing.T) {
	v := viper.GetViper()
	v.SetConfigName("default")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")
	os.Setenv("COLLECTOR_REDIS_PORT", "")
	os.Setenv("COLLECTOR_REDIS_HOST", "")
	os.Setenv("HANDLER_REDIS_PORT", "")
	os.Setenv("HANDLER_REDIS_HOST", "")
	redisConfig := CommonCfg.RedisConfig{V: v}
	assert.Equal(t, redisConfig.Port(), 6379)
	assert.Equal(t, redisConfig.Host(), "0.0.0.0")
}
