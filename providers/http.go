package providers

import (
	"context"
	"net"
	"net/http"

	"github.com/arashrasoulzadeh/go-content/routes"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewServeMux(publicRoute routes.PublicRoute, privateRoute routes.PrivateRoute) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(publicRoute.Pattern(), publicRoute)
	mux.Handle(privateRoute.Pattern(), privateRoute)
	return mux
}
