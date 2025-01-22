package components

import (
	"certifications/config"
	"certifications/internal/modules/users/http/handlers"
	"certifications/internal/modules/users/http/routers"

	"github.com/Lucasvmarangoni/logella/router"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	// "github.com/go-chi/chi/middleware"
	// "github.com/go-chi/cors"
	// "github.com/go-chi/jwtauth"
)

func Http(
	r *chi.Mux,
	userHandlers *handlers.UserHandler,
) {
	// tokenAuth := config.GetTokenAuth()
	// jwtExpiresIn := config.GetEnvInt("jwt", "expiredIn")

	router := router.NewRouter()
	userRouter := routers.NewUserRouter(userHandlers, router)

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
	r.Use(middleware.WithValue("jwt", config.GetTokenAuth()))
	// r.Use(middleware.WithValue("JwtExpiresIn", jwtExpiresIn))

	userRouter.InitializeUserRoutes(r)
}
