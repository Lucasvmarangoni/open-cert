package components

import (
	"certifications/internal/modules/users/domain/entities"
	"certifications/pkg/cache"
)

type UserCacheFactory struct {
	UserCache      *cache.Memcached[entities.User]
	CertifierCache *cache.Memcached[*entities.Certifier]
}

func NewCacheFactory() *UserCacheFactory {

	userCacheFactory := &UserCacheFactory{
		UserCache:      instances[entities.User](),
		CertifierCache: instances[*entities.Certifier](),
	}
	return userCacheFactory
}

func instances[T cache.Entity]() *cache.Memcached[T] {
	mc := cache.NewMemcached[T]("memcached-1:11211")
	return mc
}
