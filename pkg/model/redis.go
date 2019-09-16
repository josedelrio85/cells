package leads

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// Redis defines the attributes needed to perform store and retrieve operations
// for phone operations. It does implement the Storer interface.
type Redis struct {
	Pool *redis.Pool
}

// Set stores a key value on a Redis database.
//
// key: the key to store.
// value: the value to store on the given key.
//
// Returns an error if any.
func (r *Redis) Set(key, value string) error {
	redis := r.Pool.Get()
	defer redis.Close()

	if _, err := redis.Do("SET", key, value, "EX", 300); err != nil {
		return errors.Wrap(err, "error storing phone validation code")
	}

	return nil
}

// Get retrieves a value from the given key on redis.
//
// key: the key to retrieve its value.
//
// Returns the retrieved value, nil if it was not found, or an error.
func (r *Redis) Get(key string) (*string, error) {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	value, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving phone validation code")
	}

	return &value, nil
}
