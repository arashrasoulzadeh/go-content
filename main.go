package main

import (
	"net/http"

	providers2 "github.com/arashrasoulzadeh/go-content/providers"

	"github.com/arashrasoulzadeh/go-content/handlers"
	"github.com/arashrasoulzadeh/go-content/routes"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			providers2.NewMongoconnection,
			providers2.NewHTTPServer,
			providers2.NewServeMux,
			fx.Annotate(handlers.NewPublicHandler, fx.As(new(routes.PublicRoute))),
			fx.Annotate(handlers.NewPrivateHandler, fx.As(new(routes.PrivateRoute))),
			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
