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
	topic := req.URL.Query().Get("topic")
	this.appKafka.SendMessage(topic, value)

	shared.SuccessResp(w, "done!")
}
