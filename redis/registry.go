package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/citadel/citadel"
	"github.com/garyburd/redigo/redis"
)

type RedisRegistry struct {
	pool *redis.Pool
}

func NewRedisRegistry(addr, pass string) citadel.Registry {
	return &RedisRegistry{
		pool: newPool(addr, pass),
	}
}

func (r *RedisRegistry) SaveDocker(rs *citadel.Docker) error {
	key := r.getKey(rs.ID)

	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	_, err = r.do("SET", key, string(data))
	return err
}

func (r *RedisRegistry) DeleteDocker(id string) error {
	key := r.getKey(id)

	_, err := r.do("DEL", key, fmt.Sprintf("%s:reserved_cpus", key), fmt.Sprintf("%s:reserved_memory", key))
	return err
}

func (r *RedisRegistry) FetchDockers() ([]*citadel.Docker, error) {
	out := []*citadel.Docker{}

	keys, err := redis.Strings(r.do("KEYS", "citadel:resources:*"))
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		if strings.Contains(k, "reserved") {
			continue
		}

		data, err := redis.String(r.do("GET", k))
		if err != nil {
			return nil, err
		}

		var rs *citadel.Docker
		if err := json.Unmarshal([]byte(data), &rs); err != nil {
			return nil, err
		}

		out = append(out, rs)
	}
	return out, nil
}

func (r *RedisRegistry) GetTotalReservations(id string) (float64, float64, error) {
	key := r.getKey(id)

	cpus, err := redis.Float64(r.do("GET", fmt.Sprintf("%s:reserved_cpus", key)))
	if err != nil && err != redis.ErrNil {
		return 0, 0, err
	}

	memory, err := redis.Float64(r.do("GET", fmt.Sprintf("%s:reserved_memory", key)))
	if err != nil && err != redis.ErrNil {
		return 0, 0, err
	}

	return cpus, memory, nil
}

func (r *RedisRegistry) PlaceReservation(id string, c *citadel.Container) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := r.getKey(id)

	if err := conn.Send("MULTI"); err != nil {
		return err
	}

	if err := conn.Send("INCRBY", fmt.Sprintf("%s:reserved_cpus", key), c.Cpus); err != nil {
		return err
	}

	if err := conn.Send("INCRBY", fmt.Sprintf("%s:reserved_memory", key), c.Memory); err != nil {
		return err
	}

	if _, err := conn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

func (r *RedisRegistry) ReleaseReservation(id string, c *citadel.Container) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := r.getKey(id)

	if err := conn.Send("MULTI"); err != nil {
		return err
	}

	if err := conn.Send("DECRBY", fmt.Sprintf("%s:reserved_cpus", key), c.Cpus); err != nil {
		return err
	}

	if err := conn.Send("DECRBY", fmt.Sprintf("%s:reserved_memory", key), c.Memory); err != nil {
		return err
	}

	if _, err := conn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

func (r *RedisRegistry) Close() error {
	return r.pool.Close()
}

func (r *RedisRegistry) getKey(id string) string {
	return fmt.Sprintf("citadel:resources:%s", id)
}

func (r *RedisRegistry) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Do(cmd, args...)
}

func newPool(addr, pass string) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		if pass != "" {
			if _, err := c.Do("AUTH", pass); err != nil {
				return nil, err
			}
		}

		return c, nil
	}, 10)
}
