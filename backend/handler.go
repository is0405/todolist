package server

import (
	"net/http"

	"github.com/is0405/httputil"
)

type AppHandler struct {
	h func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (a AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, res, err := a.h(w, r)
	if err != nil {
		httputil.RespondErrorJson(w, status, err)
		return
	}
	httputil.RespondJSON(w, status, res)
	return
}
