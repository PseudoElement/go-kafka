package ui

import (
	"net/http"
	"os"
)

func (this *UiController) _viewRoute(w http.ResponseWriter, req *http.Request) {
	pwd, _ := os.Getwd()
	path := pwd + "/src/routes/ui/view.html"
	http.ServeFile(w, req, path)
}
