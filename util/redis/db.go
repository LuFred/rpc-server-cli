package redis

import (
	"fmt"
	"sync"
)

type DriverType int

// database alias cacher.
type _dbCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

//是否启用缓存，默认为true
var EnableCache bool

var cache = &_dbCache{cache: make(map[string]*alias)}

const (
	DriverRedis DriverType = 0
)

const (
	defaultMaxIdle        = 10
	defaultMaxActive      = 100
	defaultIdleTimeout    = 60
	defaultWait           = true
	defaultDatabase       = 0
	defaultConnectTimeout = 60
	defaultReadTimeout    = 60
	defaultWriteTimeout   = 60
)

func RegisterDb(aliasName string, driverType DriverType, host, pwd string, params ...int) error {
	setting := &connectionURL{
		Password: pwd,
		Host:     host,
	}
	if aliasName == "" {
		aliasName = "default"
	}
	switch driverType {
	case DriverRedis:
		redisDb := &RedisDB{}
		err := redisDb.Open(setting)
		if err != nil {
			go func() {
				redisDb.Close()
			}()
			return err
		}
		al, err := addAliasWithCache(aliasName, driverType, redisDb)
		if err != nil {
			go func() {
				if redisDb != nil {
					redisDb.Close()
				}
			}()
			err = fmt.Errorf("register db, %s", err.Error())
			return err
		}
		cache.add(aliasName, al)
	default:
		return ErrDriverTypeNotFound
	}
	return nil
}

//GetCache
func GetCache(aliasName string) (Database, error) {
	if aliasName == "" {
		return nil, ErrAliasNotExist
	}
	if c, ok := cache.get(aliasName); ok {
		return c.DB, nil
	} else {
		return nil, ErrAliasNotExist
	}
}

// add database alias with original name.
func (dc *_dbCache) add(name string, al *alias) (added bool) {
	dc.mux.Lock()
	defer dc.mux.Unlock()
	if _, ok := dc.cache[name]; !ok {
		dc.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dbCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dbCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}
