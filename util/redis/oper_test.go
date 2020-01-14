package redis

import "testing"

type TestDemo struct {
	Age   int32      `redis:"age"`
	Name  string     `redis:"name"`
	Skill []TestDemo `redis:"skill"`
}

func init() {
	TestOpen(&testing.T{})
}
func TestOpen(t *testing.T) {
	err := RegisterDb("", DriverRedis, "127.0.0.1:6380", "")
	if err != nil {
		t.Errorf("open err =%s", err.Error())
	} else {
		t.Log("open success")
	}
}
func TestSSet(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	demo := &TestDemo{
		Age:  3,
		Name: "xxx",
	}
	err = db.SSet("test", "key1", demo, 30)
	if err != nil {
		t.Errorf("sset err =%s", err.Error())
	} else {
		t.Log("sset success")
	}
}

func TestExists(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	b, err := db.Exists("test", "key1")
	if err != nil {
		t.Errorf("sset err =%s", err.Error())
	} else {
		t.Logf("exist success %t", b)
	}
}

func TestSGet(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	b, err := db.SGet("test", "key1")
	if err != nil {
		t.Errorf("sset err =%s", err.Error())
	} else {
		t.Logf("sget success %s", b)
	}
}

func TestDel(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	b, err := db.Del("test", "key1")
	if err != nil {
		t.Errorf("sset err =%s", err.Error())
	} else {
		t.Logf("del success %d", b)
	}
}

func TestHSet(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	demo := &TestDemo{
		Age:  3,
		Name: "33",
		Skill: []TestDemo{
			TestDemo{Age: 10, Name: "cn1"},
			TestDemo{Age: 20, Name: "cn2"}},
	}
	b, err := db.HSet("test", "hkey1", "field1", demo)

	if err != nil {
		t.Errorf("hset err =%s", err.Error())
	} else {
		t.Logf("hset success %d", b)
	}
}

func TestExpire(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	r, err := db.Expire("test", "hkey1", 30)
	if err != nil {
		t.Errorf("hset err =%s", err.Error())
	} else {
		t.Logf("expire success  %d", r)
	}
}

func TestHGet(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}

	b, err := db.HGetVals("test", "hkey1")

	if err != nil {
		t.Errorf("hset err =%s", err.Error())
	} else {
		t.Logf("hset success %v", b[0])
	}
}

func TestHGetAll(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	var r map[string]string
	r, err = db.HGetAll("test", "demo")

	if err != nil {
		t.Errorf("TestHGetAll err =%s", err.Error())
	} else {
		t.Logf("TestHGetAll success %v", r)
	}
}

func TestHExist(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	r, err := db.HExist("test", "demo", "name")

	if err != nil {
		t.Errorf("TestHExist err =%s", err.Error())
	} else {
		t.Logf("TestHExist success %t", r)
	}
}

func TestHDel(t *testing.T) {
	db, err := GetCache("default")
	if err != nil {
		t.Error(err)
	}
	r, err := db.HDel("test", "demo", "name", "name1", "name2")

	if err != nil {
		t.Errorf("TestHDel err =%s", err.Error())
	} else {
		t.Logf("TestHDel success %d", r)
	}
}
