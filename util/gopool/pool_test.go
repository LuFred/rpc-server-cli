package pool

import (
	"testing"
)

type User struct {
	Name string
	Age  int32
}

func TestGo(t *testing.T) {
	InitGoPool(10)
	u := User{Name: "xx", Age: 3}
	Go(func() {
		u.Name = "ssdasd"
		panic("hahaha")
		t.Logf("user %s", u.Name)
	})

}
