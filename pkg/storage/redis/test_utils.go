package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sciclon2/kafka-lag-go/pkg/config"
	"github.com/stretchr/testify/mock"
)

// defaultRedisConfig returns a standard Redis configuration.
func defaultRedisConfig() *config.Config {
	return &config.Config{
		Storage: struct {
			Type  string `yaml:"type"`
			Redis struct {
				Address              string `yaml:"address"`
				Port                 int    `yaml:"port"`
				ClientRequestTimeout string `yaml:"client_request_timeout"`
				ClientIdleTimeout    string `yaml:"client_idle_timeout"`
				RetentionTTLSeconds  int    `yaml:"retention_ttl_seconds"`
			} `yaml:"redis"`
		}{
			Type: "redis",
			Redis: struct {
				Address              string `yaml:"address"`
				Port                 int    `yaml:"port"`
				ClientRequestTimeout string `yaml:"client_request_timeout"`
				ClientIdleTimeout    string `yaml:"client_idle_timeout"`
				RetentionTTLSeconds  int    `yaml:"retention_ttl_seconds"`
			}{
				Address:              "localhost",
				Port:                 6379,
				ClientRequestTimeout: "5s",
				ClientIdleTimeout:    "5m",
			},
		},
	}
}

// MockRedisClient is a mock implementation of the RedisClient interface.
type MockRedisClient struct {
	mock.Mock
}
type MockPipeliner struct {
	mock.Mock
	redis.Pipeliner // Embed the Pipeliner interface to satisfy the interface requirements
}

func (m *MockPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]redis.Cmder), args.Error(1)
}

func (m *MockPipeliner) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	// Properly initialize the Cmd object
	cmd := redis.NewCmd(ctx)
	m.Called(ctx, sha1, keys, args)
	return cmd
}

// ZRangeWithScores is a mock implementation of the ZRangeWithScores method from the Pipeliner interface.
func (m *MockPipeliner) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).(*redis.ZSliceCmd)
}

type RedisManagerWithMock struct {
	*RedisManager
	mock.Mock
}

func (m *MockRedisClient) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	args := m.Called(ctx, script)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	allArgs := append([]interface{}{ctx, sha1, keys}, args...)
	argsCalled := m.Called(allArgs...)
	return argsCalled.Get(0).(*redis.Cmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) Pipeline() redis.Pipeliner {
	args := m.Called()
	return args.Get(0).(redis.Pipeliner)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

// createMockRedisClient returns a new MockRedisClient instance.
func createMockRedisClient() *MockRedisClient {
	return new(MockRedisClient)
}
