package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
	CommonCfg "github.com/webalytic.go/common/config"
)

func SetUp() *viper.Viper {
	viperObj := viper.GetViper()
	viperObj.SetConfigName("default")
	viperObj.SetConfigType("yaml")
	viperObj.AddConfigPath("./common/config")
	viperObj.AddConfigPath("./")
	err := viperObj.ReadInConfig()
	if err != nil {
		fmt.Println("Unable to read config")
	}
	viperObj.SetDefault("APP_PREFIX", "app")
	viperObj.SetEnvPrefix(viperObj.GetString("APP_PREFIX"))
	return viperObj
}

func TestLoadRedisConfigFromDefaultConfig(t *testing.T) {
	viperObj := SetUp()
	redisConfig := CommonCfg.RedisConfig{Viper: viperObj}
	assert.Equal(t, redisConfig.Port(), 6379)
	assert.Equal(t, redisConfig.Host(), "0.0.0.0")
}

func TestLoadRedisEnvShouldBePriorForCollector(t *testing.T) {
	viperObj := SetUp()
	redisConfig := CommonCfg.RedisConfig{Viper: viperObj}
	viperObj.SetEnvPrefix("REDIS")
	os.Setenv("REDIS_PORT", "1234")
	os.Setenv("REDIS_HOST", "127.0.0.1")
	assert.Equal(t, 1234, redisConfig.Port())
	assert.Equal(t, "127.0.0.1", redisConfig.Host())
}

func TestLoadRedisEnvShouldBePriorForHandler(t *testing.T) {
	viperObj := SetUp()
	redisConfig := CommonCfg.RedisConfig{Viper: viperObj}
	viperObj.SetEnvPrefix("REDIS")
	os.Setenv("REDIS_PORT", "1234")
	os.Setenv("REDIS_HOST", "127.0.0.1")
	assert.Equal(t, 1234, redisConfig.Port())
	assert.Equal(t, "127.0.0.1", redisConfig.Host())
}

func TestNoDefaultFileFound(t *testing.T) {
	viperObj := SetUp()
	viperObj.SetConfigName("default")
	viperObj.SetConfigType("yaml")
	viperObj.AddConfigPath("./")
	os.Setenv("REDIS_PORT", "")
	os.Setenv("REDIS_HOST", "")
	os.Setenv("REDIS_PORT", "")
	os.Setenv("REDIS_HOST", "")
	redisConfig := CommonCfg.RedisConfig{Viper: viperObj}
	assert.Equal(t, redisConfig.Port(), 6379)
	assert.Equal(t, redisConfig.Host(), "0.0.0.0")
}
