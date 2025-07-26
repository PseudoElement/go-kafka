package gateway

import (
	"net/http"

	"github.com/pseudoelement/go-kafka/src/shared"
)

func (this *GatewayController) _testRoute(w http.ResponseWriter, req *http.Request) {
	resp := ValueResponse{Value: "Chi router works!"}
	shared.SuccessResp(w, resp)
}

func (this *GatewayController) _requestRoute(w http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("value")
	resp := ValueResponse{Value: "Kafka handled " + value}

	this.appKafka.SendMessage("sintol-topic", resp)
	shared.SuccessResp(w, resp)
}
