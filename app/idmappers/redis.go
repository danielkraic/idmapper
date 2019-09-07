package idmappers

import (
	"fmt"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/go-redis/redis"
)

// NewRedisIDMapper creates IDMapper that reads data from redis
func NewRedisIDMapper(client *redis.Client, hashName string) (*idmapper.IDMapper, error) {
	return idmapper.NewIDMapper(&redisSource{client: client, hashName: hashName})
}

type redisSource struct {
	client   *redis.Client
	hashName string
}

func (r *redisSource) Read() (idmapper.ValuesMap, error) {
	result, err := r.client.HGetAll(r.hashName).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to HGET hash %s: %s", r.hashName, err)
	}

	return result, nil
}
