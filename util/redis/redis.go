package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisDB struct {
	db *redis.Pool
}

type HVal struct {
	Field string
	Value string
}

func (cache *RedisDB) Open(settings *connectionURL) error {
	cache.db = &redis.Pool{
		MaxIdle:     defaultMaxIdle,
		MaxActive:   defaultMaxActive,
		IdleTimeout: defaultIdleTimeout * time.Second,
		Wait:        defaultWait,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp",
				settings.Host,
				redis.DialPassword(settings.Password),
				redis.DialDatabase(settings.Database),
				redis.DialConnectTimeout(defaultConnectTimeout*time.Second),
				redis.DialReadTimeout(defaultReadTimeout*time.Second),
				redis.DialWriteTimeout(defaultWriteTimeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	return nil
}
func (cache *RedisDB) Close() error {
	if cache.db != nil {
		return cache.db.Close()
	} else {
		return ErrDbNotFound
	}
}
func (cache *RedisDB) Exists(Prefix, key string) (bool, error) {
	var res bool
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	res, err := redis.Bool(con.Do("EXISTS", Prefix+key))
	if err != nil {
		return res, err
	}
	return res, nil
}

func (cache *RedisDB) Expire(Prefix, key string, expireSecond int) (int64, error) {
	var res int64
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	res, err := redis.Int64(con.Do("EXPIRE", Prefix+key, expireSecond))
	if err != nil {
		return res, err
	}
	return res, nil
}
func (cache *RedisDB) SGet(Prefix, key string) (string, error) {
	var res string
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	res, err := redis.String(con.Do("GET", Prefix+key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}
	return res, nil
}
func (cache *RedisDB) SSet(Prefix, key string, value interface{}, expireSecond int) error {
	if key == "" {
		return ErrKeyNotBeenNil
	}
	if value == nil {
		return ErrValueNotBeenNil
	}
	v, err := convertInterfaceToString(value)
	if err != nil {
		return err
	}
	if v == "" {
		return ErrValueNotBeenNil
	}
	con := cache.db.Get()
	if con.Err() != nil {
		return con.Err()
	}
	defer con.Close()
	if expireSecond > 0 {
		_, err = con.Do("SET", Prefix+key, v, "EX", expireSecond)
	} else {
		_, err = con.Do("SET", Prefix+key, v)
	}
	if err != nil {
		return err
	}
	return nil
}
func (cache *RedisDB) Del(Prefix, key string) (int64, error) {
	var res int64
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	var err error
	index, err := con.Do("DEL", Prefix+key)
	if err != nil {
		return 0, err
	}
	return index.(int64), nil
}
func (cache *RedisDB) HSet(Prefix, key, field string, value interface{}) (int64, error) {
	var res int64
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	var err error
	v, err := convertInterfaceToString(value)
	if err != nil {
		return res, err
	}
	if v == "" {
		return res, ErrValueNotBeenNil
	}
	index, err := con.Do("HSET", Prefix+key, field, v)
	if err != nil {
		return 0, err
	}
	return index.(int64), nil
}

// hash bucket:set ,with expiration
func (cache *RedisDB) HSetExpire(Prefix, key, field, data string, expireSecond int) (int64, error) {
	var res int64
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	if expireSecond > 0 {
		index, err := con.Do("HSET", Prefix+key, field, data, "EX", expireSecond)
		if err != nil {
			return 0, err
		}
		return index.(int64), nil
	} else {
		index, err := con.Do("HSET", Prefix+key, field, data)
		if err != nil {
			return 0, err
		}
		return index.(int64), nil
	}
}

//todo 不支持负责结构体，带研究
// func (cache *RedisDB)HMSet(Prefix,key,value interface{}) (error){
//
// 	con:=cache.db.Get()
// 	if con.Err()!=nil{
// 		return con.Err()
// 	}
// 	defer con.Close()
// 	var err error
// 	v,err:=convertInterfaceToString(value)
// 	if err!=nil{
// 		return  err
// 	}
// 	if v==""{
// 		return ErrValueNotBeenNil
// 	}
// 	_, err = con.Do("HMSET",redis.Args{}.Add(fmt.Sprintf("%s%s",Prefix,key)).AddFlat(value)...)
// 	if err != nil {
// 		return err
// 	}
// 	return  nil
// }

func (cache *RedisDB) HGetVals(Prefix, key string) ([]string, error) {
	var res []string
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	var err error
	vs, err := redis.Values(con.Do("HVALS", Prefix+key))
	if err != nil && err != redis.ErrNil {
		return res, err
	}
	c := len(vs)
	res = make([]string, c)
	for i, v := range vs {
		res[i] = string(v.([]byte)[:])
	}
	return res, nil
}
func (cache *RedisDB) HGetAll(Prefix string, key string) (map[string]string, error) {
	var res map[string]string
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	res, err := redis.StringMap(con.Do("HGETALL", Prefix+key))
	if err != nil {
		return res, fmt.Errorf("Cache: StringMap err=%s", err.Error())
	}
	return res, nil
}
func (cache *RedisDB) HDel(Prefix, key string, fields ...string) (int64, error) {
	var res int64
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	var err error
	count := len(fields)
	if count < 1 {
		return 0, nil
	}
	list := make([]interface{}, count+1)
	list[0] = Prefix + key
	for i, v := range fields {
		list[i+1] = v
	}
	index, err := con.Do("HDEL", list...)
	if err != nil {
		return 0, err
	}
	return index.(int64), nil
}
func (cache *RedisDB) HExist(Prefix, key, field string) (bool, error) {
	var res bool
	con := cache.db.Get()
	if con.Err() != nil {
		return res, con.Err()
	}
	defer con.Close()
	res, err := redis.Bool(con.Do("HEXISTS", Prefix+key, field))
	if err != nil {
		return res, err
	}
	return res, nil
}

// hash bucket: get all values,each with field,value
func (cache *RedisDB) HGetValsByMulKey(Prefix string, keys ...string) (hv []HVal, err error) {
	con := cache.db.Get()
	if con.Err() != nil {
		return nil, con.Err()
	}
	defer con.Close()
	for _, key := range keys {
		result, err := redis.Values(con.Do("HGETALL", Prefix+key))
		if err != nil {
			continue
		}
		for k, v := range result {
			if k%2 == 1 {
				hv = append(hv, HVal{
					Field: fmt.Sprintf("%s", result[k-1]),
					Value: fmt.Sprintf("%s", v),
				})
			}
		}

	}

	return
}

//todo value转结构体，待完成
// var structSpecCache = make(map[string]string)
//
// func scanStruct(src []interface{},dest interface{})error{
// 	d := reflect.ValueOf(dest)
// 	if d.Kind() != reflect.Ptr || d.IsNil() {
// 		return ErrScanStructValue
// 	}
// 	dst := d.Elem()
// 	if len(src)%2 != 0 {
// 		return fmt.Errorf("Cache: number of values not a multiple of 2")
// 	}
// 	//遍历获取所有一级字段
// 	for i := 0; i < len(src); i += 2 {
// 		val := src[i+1]
// 		if val== nil {
// 			continue
// 		}
// 		key, ok := src[i].([]byte)
// 		if !ok {
// 			return fmt.Errorf("Cache: key %d not a bulk string value", i)
// 		}
// 		structSpecCache[string(key)]=string(val.([]byte))
// 	}
// 	//遍历目标对象
// 	types:=dst.Type()
// 	for i := 0; i < types.NumField(); i++ {
// 		field := types.Field(i)
// 		fValue:=dst.FieldByName(field.Name)
// 		fs,err:=compileField(field)
// 		if err!=nil{
// 			return fmt.Errorf("Cache: key %d not a bulk string value", i)
// 		}
// 		refVal,err:=bind(fs.key,structSpecCache[fs.key],dst.FieldByName(field.Name).Type())
// 		if err!=nil{
// 			return fmt.Errorf("Cache: bind value error : %s",err.Error())
// 		}
// 		fValue.Set(refVal)
// 	}
// 	return nil
// }
// func bind(key string ,val string,typ reflect.Type) (reflect.Value,error){
// 	rv := reflect.Zero(typ)
// 	switch typ.Kind() {
// 	case reflect.Ptr:
// 		vs, err := bind(key,val,typ.Elem())
// 		if err != nil {
// 			return reflect.Zero(typ), err
// 		}
// 		return vs.Addr(), nil
// 	case reflect.Slice:
// 		sl,err:=slice(val)
// 		if err!=nil{
// 			return reflect.Zero(typ), err
// 		}
//
// 		return bindSlice(key,sl,typ)
// 	case reflect.Struct:
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return bindInt(val,typ),nil
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		return bindUint(val,typ),nil
// 	case reflect.Float32, reflect.Float64:
// 		return bindFloat(val,typ),nil
// 	case reflect.String:
// 		return bindString(val,typ),nil
// 	case reflect.Bool:
// 		return bindBool(val,typ),nil
// 	}
// 	return rv, nil
// }
// type Value struct {
// 	data   interface{}
// 	exists bool
// }
// func slice(vals string)([]*Value, error) {
// 	var slice []*Value
// 	fmt.Print("======",vals)
// 	vr:=bytes.NewReader([]byte(vals))
// 	var data interface{}
// 	err := json.NewDecoder(vr).Decode(data)
// 	if err != nil {
//
// 		return slice, err
// 	}
//
// 	var valid bool
// 	switch data.(type) {
// 	case []interface{}:
// 		valid=true
// 	}
// 	if !valid{
// 		return slice,errors.New("bind slice Not an array")
// 	}
// 	for _, item := range data.([]interface{}) {
// 		val := Value{item, true}
// 		slice = append(slice, &val)
// 	}
// 	return slice,nil
// }
// func bindSlice(key string,vals []*Value,  typ reflect.Type) (reflect.Value, error) {
// 	sv := reflect.MakeSlice(typ, 0, len(vals))
// 	sliceValueType := typ.Elem()
// 	for _, v := range vals {
// 		childValue, err := bind(key,v.data.(string),sliceValueType)
// 		if err != nil {
// 			return reflect.Zero(typ), err
// 		}
// 		sv = reflect.Append(sv, childValue)
// 	}
// 	return sv, nil
// }
// func bindInt(val string, typ reflect.Type) reflect.Value {
// 	iv, err := strconv.ParseInt(val, 10, 64)
// 	if err != nil {
// 		return reflect.Zero(typ)
// 	}
// 	pValue := reflect.New(typ)
// 	pValue.Elem().SetInt(iv)
// 	return pValue.Elem()
// }
//
// func bindUint(val string, typ reflect.Type) reflect.Value {
// 	uv, err := strconv.ParseUint(val, 10, 64)
// 	if err != nil {
// 		return reflect.Zero(typ)
// 	}
// 	pValue := reflect.New(typ)
// 	pValue.Elem().SetUint(uv)
// 	return pValue.Elem()
// }
//
// func bindFloat(val string, typ reflect.Type) reflect.Value {
// 	fv, err := strconv.ParseFloat(val, 64)
// 	if err != nil {
// 		return reflect.Zero(typ)
// 	}
// 	pValue := reflect.New(typ)
// 	pValue.Elem().SetFloat(fv)
// 	return pValue.Elem()
// }
// func bindBool(val string, typ reflect.Type) reflect.Value {
// 	val = strings.TrimSpace(strings.ToLower(val))
// 	switch val {
// 	case "true", "on", "1":
// 		return reflect.ValueOf(true)
// 	}
// 	return reflect.ValueOf(false)
// }
//
// func bindString(val string, typ reflect.Type) reflect.Value {
// 	return reflect.ValueOf(val)
// }
//
// type fieldSpec struct {
// 	key      string
// 	index     []int
// 	omitEmpty bool
// 	value string
// }
//
// type structSpec struct {
// 	m map[string]*fieldSpec
// }
//
// func compileField(field reflect.StructField)(fs *fieldSpec,err error) {
// 		fs=new(fieldSpec)
// 		switch  {
// 		case field.PkgPath!=""&&!field.Anonymous:
// 			// Ignore unexported fields.
// 			return
// 		case field.Anonymous:
// 			if field.Type.Kind()==reflect.Struct{
// 			}
// 		default:
// 			fs.key= field.Name
// 			tag:=field.Tag.Get("redis")
// 			if tag==""{
// 				return
// 			}
// 			p:=strings.Split(tag,",")
// 			if len(p) > 0 {
// 				if p[0] == "-" {
// 					return
// 				}
// 				if len(p[0]) > 0 {
// 					fs.key = p[0]
// 				}
// 				for _, s := range p[1:] {
// 					switch s {
// 					case "omitempty":
// 						fs.omitEmpty = true
// 					default:
//
// 					}
// 				}
// 			}
// 		}
// 		return
// }
//
//
//
