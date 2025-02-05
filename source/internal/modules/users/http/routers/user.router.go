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
	UserHandler *handlers.UsersHandler
	Router      *router.Router
}

func NewUserRouter(
	userHandler *handlers.UsersHandler,
	router *router.Router,
) *UserRouter {
	return &UserRouter{
		UserHandler: userHandler,
		Router:      router,
	}
}

func (u *UserRouter) InitializeUserRoutes(r chi.Router) {
	prefix := "/authn"
	r.Route(prefix, func(r chi.Router) {
		u.Router.Method(router.POST).Prefix(prefix).InitializeRoute(r, "/create", u.UserHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				10,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Router.Method("POST").Prefix(prefix).InitializeRoute(r, "/jwt", u.UserHandler.Authentication)
		})

	})
}

func (u *UserRouter) UserTOTPVerifyRoute(r chi.Router) {
	u.Router.Method("POST").Prefix("/2fa").InitializeRoute(r, "/verify", u.UserHandler.TwoFactorVerify)
}

func (u *UserRouter) UserRoutes(r chi.Router) {

	prefix := "/totp"
	r.Route(prefix, func(r chi.Router) {
		u.Router.Method("GET").Prefix(prefix).InitializeRoute(r, "/generate", u.UserHandler.TwoFactorAuthn)
		u.Router.Method("PATCH").Prefix(prefix).InitializeRoute(r, "/disable", u.UserHandler.TwoFactorDisable)
	})

	prefix = "/user"
	r.Route(prefix, func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				5,
				30*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Router.Method("GET").Prefix(prefix).InitializeRoute(r, "/me", u.UserHandler.Me)
			u.Router.Method("PUT").Prefix(prefix).InitializeRoute(r, "/update", u.UserHandler.Update)
		})

		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				4,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByRealIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Router.Method("DELETE").Prefix(prefix).InitializeRoute(r, "/del", u.UserHandler.Delete)
		})
	})
}

func (u *UserRouter) CertifierRoute(r chi.Router) {
	prefix := "/certifier"
	r.Route(prefix, func(r chi.Router) {
		u.Router.Method(router.POST).Prefix(prefix).InitializeRoute(r, "/create", u.UserHandler.CreateCertifier)
	})
}
