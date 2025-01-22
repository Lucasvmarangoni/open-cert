package routers

import (
	"certifications/internal/modules/users/http/handlers"
	"net/http"
	"time"

	"github.com/Lucasvmarangoni/logella/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
)

type UserRouter struct {
	UserHandler    *handlers.UserHandler
	Router         *router.Router
}

func NewUserRouter(
	userHandler *handlers.UserHandler,
	router *router.Router,
) *UserRouter {
	return &UserRouter{
		UserHandler:    userHandler,
		Router:         router,
	}
}

func (u *UserRouter) InitializeUserRoutes(r chi.Router) {
	prefix := "/authn"
	r.Route(prefix, func(r chi.Router) {
		u.Router.Method(router.POST).Prefix(prefix).InitializeRoute(r, "/create", u.UserHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				2,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			// u.Router.Method("POST").Prefix(prefix).InitializeRoute(r, "/jwt", u.UserHandler.Authentication)
		})

	})
}
