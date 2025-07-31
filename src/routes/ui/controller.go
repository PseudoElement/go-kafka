package ui

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/pseudoelement/go-kafka/src/kafka"
)

type UiController struct {
	r        *chi.Mux
	appCtx   context.Context
	appKafka *kafka.AppKafka
}

func NewUiController(r *chi.Mux, appKafka *kafka.AppKafka, appCtx context.Context) *UiController {
	uiController := &UiController{r: r, appKafka: appKafka, appCtx: appCtx}
	uiController.appKafka.AddListener("ui_controller")
	go uiController.listenChan()

	return uiController
}

func (this *UiController) SetRoutes() {
	this.r.Route("/ui", func(r chi.Router) {
		r.With().Get("/page", this._viewRoute)
	})
}

func (this *UiController) listenChan() {
	for msg := range this.appKafka.ConsumerChan("ui_controller") {
		log.Printf("%s UiController message - %s.\n", msg.Topic, string(msg.Value))
	}
}
