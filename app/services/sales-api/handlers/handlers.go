package handlers

import (
	"expvar"
	"github.com/dtherhtun/service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/dtherhtun/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/dtherhtun/service/business/web/mid"
	"github.com/dtherhtun/service/foundation/web"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *web.App {

	// Construct the web.app which holds all routes.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
	)

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

// v1 binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
}
