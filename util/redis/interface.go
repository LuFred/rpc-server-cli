package redis

type Database interface {
	Open(settings *connectionURL) error
	Close() error
	//判断key是否存在
	Exists(Prefix, key string) (bool, error)
	Expire(Prefix, key string, expireSecond int) (int64, error)
	SGet(Prefix, key string) (string, error)
	SSet(Prefix, key string, value interface{}, expireSecond int) error
	Del(Prefix, key string) (int64, error)
	HSet(Prefix, key, field string, value interface{}) (int64, error)
	HSetExpire(Prefix, key, field, data string, expireSecond int) (int64, error)

	//HMSet(Prefix,key,value interface{}) ( error)

	//HGetVals Redis Hvals 命令返回哈希表所有域(field)的值。
	HGetVals(Prefix, key string) ([]string, error)
	HGetAll(Prefix string, key string) (map[string]string, error)
	HGetValsByMulKey(Prefix string, keys ...string) (hv []HVal, err error)
	HDel(Prefix, key string, fields ...string) (int64, error)
	HExist(Prefix, key string, fields string) (bool, error)
}
