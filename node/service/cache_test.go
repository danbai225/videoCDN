package service

import (
	"testing"
)

func TestCache(t *testing.T) {
	err := Cache("123", []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	cache, err := GetCache("123")
	if err != nil {
		t.Fatal(err)
	}
	if string(cache) != "test" {
		t.Fatal("缓存错误")
	}
}
