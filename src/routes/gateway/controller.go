package gateway

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pseudoelement/go-kafka/src/kafka"
)

type GatewayController struct {
	r        *chi.Mux
	appCtx   context.Context
	appKafka *kafka.AppKafka
}

func NewGatewayController(r *chi.Mux, appKafka *kafka.AppKafka, appCtx context.Context) *GatewayController {
	return &GatewayController{r: r, appCtx: appCtx}
}

func (this *GatewayController) SetRoutes() {
	this.r.Route("/gateway", func(r chi.Router) {
		r.With().Get("/test", this._testRoute)
		r.With().Post("/request", this._requestRoute)
	})
}
