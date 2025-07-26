package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/pseudoelement/go-kafka/src/kafka"
	"github.com/pseudoelement/go-kafka/src/middlewares"
	"github.com/pseudoelement/go-kafka/src/routes/gateway"
	"github.com/pseudoelement/go-kafka/src/shared"
)

func _notFoundRoute(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fullUrl := url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   r.URL.RequestURI(),
	}

	msg := fullUrl.String() + " route not found!"
	shared.FailedResp(w, msg, http.StatusNotFound)
}

func stopKafka(cancel context.CancelFunc, ctx context.Context) {
	secString := os.Args[1]
	sec, err := strconv.Atoi(secString)
	if err != nil {
		panic(secString + " is invalid duration in secs value")
	}

	time.Sleep(time.Second * time.Duration(sec))
	cancel()

	println("Kafka was stopped.")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	appCtx, cancel := context.WithCancel(context.Background())
	appKafka := kafka.NewAppKafka(appCtx)

	go stopKafka(cancel, appCtx)

	router := chi.NewRouter()
	apiRouterV1 := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json", "text/xml", "text/plain"))
	router.Use(middleware.CleanPath)
	router.Use(middlewares.XApiTokenMiddleware)

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins:   []string{"https://your-production-domain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-api-token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.NotFound(_notFoundRoute)

	gatewayController := gateway.NewGatewayController(apiRouterV1, appKafka, appCtx)

	gatewayController.SetRoutes()

	router.Mount("/api/v1", apiRouterV1)

	println("Start server...")

	http.ListenAndServe(":8080", router)
}
