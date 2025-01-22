package cache

import (
	"certifications/internal/modules/users/domain/entities"
	"encoding/json"

	"github.com/Lucasvmarangoni/logella/err"
	"github.com/bradfitz/gomemcache/memcache"
)

type Entity interface {
	*entities.User | *entities.Certifier | *entities.Reputation | *entities.Experiences | *entities.Qualification
}

type Mencacher[T Entity] interface {
	Set(key string, i T) error
	Get(key string) (*memcache.Item, error)
	GetMulti(key []string) (map[string]*memcache.Item, error)
	Delete(key string) error
	SetUnique(key string) error
}

type Memcached[T Entity] struct {
	Client *memcache.Client
}

func NewMemcached[T Entity](server ...string) *Memcached[T] {
	client := memcache.New(server...)
	return &Memcached[T]{
		Client: client,
	}
}

func (m *Memcached[T]) Set(key string, i T) error {
	value, err := json.Marshal(i)
	if err != nil {
		return errs.Ctx(err, "Error marshaling data to JSON")
	}
	item := &memcache.Item{
		Key:   key,
		Value: value,
	}
	err = m.Client.Set(item)
	if err != nil {
		return errs.Ctx(err, "Error setting cache")
	}
	return nil
}

func (m *Memcached[bool]) SetUnique(key string) error {

	item := &memcache.Item{
		Key: key,
	}

	err := m.Client.Set(item)
	if err != nil {
		return errs.Ctx(err, "Error setting cache")
	}
	return nil
}

func (m *Memcached[T]) Get(key string) (*memcache.Item, error) {
	item, err := m.Client.Get(key)

	if err != nil {
		return nil, errs.Ctx(err, "Error get cache")
	}
	return item, nil
}

func (m *Memcached[T]) GetUnique(key string) error {
	_, err := m.Client.Get(key)

	if err != nil {
		return errs.Ctx(err, "Error get cache")
	}
	return nil
}

func (m *Memcached[T]) GetMulti(key []string) (map[string]*memcache.Item, error) {

	items, err := m.Client.GetMulti(key)
	if err != nil {
		return nil, errs.Ctx(err, "Error get multi cache")
	}
	return items, nil
}

func (m *Memcached[T]) Delete(key string) error {
	err := m.Client.Delete(key)
	if err != nil {
		return errs.Ctx(err, "Error to delete cache")
	}
	return nil
}
