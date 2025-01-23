package services

import (
	"certifications/internal/modules/users/domain/entities"
	// pkg_entities "certifications/pkg/entities"
	"context"

	"github.com/Lucasvmarangoni/logella/err"
)

// func (u *UserService) FindById(id string, ctx context.Context) (*entities.User, error) {

// 	u.cacheUser = u.getCache(id)
// 	if cachedUser, ok := u.returnCachedUserIfExists(); ok {
// 		err := u.decrypt(cachedUser)
// 		if err != nil {
// 			return nil, errs.Ctx(err, "u.decrypt")
// 		}
// 		return cachedUser, nil
// 	}

// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	parsedId, err := pkg_entities.ParseID(id)
// 	if err != nil {
// 		return nil, errs.Ctx(err, "pkg_entities.ParseID")
// 	}
// 	user, err := u.Repository.FindById(parsedId, ctx)
// 	if err != nil {
// 		return nil, errs.Ctx(err, "Repository.FindById")
// 	}

// 	u.setToMemcacheIfNotNil(user)

// 	err = u.decrypt(user)
// 	if err != nil {
// 		return nil, errs.Ctx(err, "u.decrypt")
// 	}

// 	return user, nil
// }

func (u *UserService) FindByEmail(email string, ctx context.Context) (*entities.User, error) {

	// if cachedUser, ok := u.returnCachedUserIfExists(); ok {
	// 	return cachedUser, nil
	// }

	if ctx == nil {
		ctx = context.Background()
	}

	user, err := u.Repository.FindByEmail(email, ctx)
	if err != nil {
		return nil, errs.Ctx(err, "Repository.FindByEmail")
	}
	u.setToMemcacheIfNotNil(user)
	return user, nil
}
