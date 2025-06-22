package cache

import (
	"github.com/spf13/viper"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"sync"
)

type keyType interface {
	string | int | bool
}

type Key[T keyType] struct {
	key string
}

var cache = viper.New()
var cacheLock sync.Mutex

func NewKey[T keyType](key string) *Key[T] {
	return &Key[T]{key: key}
}

func (k *Key[T]) Save(v T) {
	cache.Set(k.key, v)
	err := save()
	if err != nil {
		log.Error("set cache failed", "err", err)
	}
}

func (k *Key[T]) Get() T {
	v := cache.Get(k.key)
	if v == nil {
		return *new(T)
	}
	return v.(T)
}

func save() error {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err := cache.WriteConfig()
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func Init(path string) {
	if !fileutil.Exist(path) {
		err := fileutil.Write(path, []byte("{}"))
		errutil.Check(err)
	}

	cache.SetConfigType("json")
	cache.SetConfigFile(path)
	err := viper.ReadInConfig()
	errutil.Check(err)
}
