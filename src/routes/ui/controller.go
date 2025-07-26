package ui

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pseudoelement/go-kafka/src/kafka"
)

type UiController struct {
	r      *chi.Mux
	appCtx context.Context
}

func NewUiController(r *chi.Mux, appKafka *kafka.AppKafka, appCtx context.Context) *UiController {
	return &UiController{r: r, appCtx: appCtx}
}

func (this *UiController) SetRoutes() {
	this.r.Route("/ui", func(r chi.Router) {
		r.With().Get("/page", this._viewRoute)
	})
}
