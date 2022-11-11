package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	// api server
	fx.Provide(NewUserRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewApiRoutes),
	fx.Provide(NewPprofRoutes),

	// file server
	fx.Provide(NewFileRoutes),
)

// ApiRoutes contains multiple routes
type ApiRoutes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewApiRoutes(
	userRoutes *UserRoutes,
	authRoutes *AuthRoutes,
	pprofRoutes *PprofRoutes,
) ApiRoutes {
	return ApiRoutes{
		userRoutes,
		authRoutes,
		pprofRoutes,
	}
}

// Setup all the route
func (r ApiRoutes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
