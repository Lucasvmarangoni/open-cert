package components

import (
	"certifications/config"
	"certifications/internal/middlewares"
	"certifications/internal/modules/users/domain/services"
	"certifications/internal/modules/users/http/handlers"
	"certifications/internal/modules/users/http/routers"

	"github.com/Lucasvmarangoni/logella/router"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

func Http(
	r *chi.Mux,
	userService *services.UserService,
	certifierService *services.CertifierService,
) {
	tokenAuth := config.GetTokenAuth()
	tokenAuthTemp := config.GetTokenAuthTemp()
	redisPassword := config.GetEnvString("password", "redis")

	userHandler := handlers.NewUserHandler(userService, certifierService)

	router := router.NewRouter()
	userRouter := routers.NewUserRouter(userHandler, router)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// r.Use(middlewares.WAF())
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)	
	r.Use(middlewares.RequestIDMiddleware)

	userRouter.InitializeUserRoutes(r)

	r.Route("/2fa", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuthTemp))
		r.Use(jwtauth.Authenticator)
		// r.Use(mw.Authorizer())
		r.Use(middlewares.NewUserRateLimit("redis:6379", redisPassword).Handler())
		userRouter.UserTOTPVerifyRoute(r)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		// r.Use(mw.Authorizer())
		r.Use(middlewares.NewUserRateLimit("redis:6379", redisPassword).Handler())
		userRouter.UserRoutes(r)
		userRouter.CertifierRoute(r)
	})
}
