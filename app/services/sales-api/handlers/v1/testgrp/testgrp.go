// Package testgrp contains all the test handlers.
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"go.uber.org/zap"

	"github.com/dtherhtun/service/business/sys/validate"
	"github.com/dtherhtun/service/foundation/web"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		//return errors.New("untrusted error")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		//return web.NewShutdownError("restart service")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
