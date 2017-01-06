package handlers

import (
	"app/interfaces/errs"
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"

	"io"

	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
)

// decodeReq decodes request's body to given interface
func decodeReq(r *http.Request, to interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(to); err != nil {
		if err != io.EOF {
			return errs.Wrap(err)
		}
	}
	return nil
}

type response struct {
	Result interface{} `json:"result"`
}

func qParam(k string, r *http.Request) string {
	values := r.URL.Query()[k]

	if len(values) != 0 {
		return values[0]
	}

	return ""
}

func muxVarMustInt(k string, r *http.Request) int {
	i, err := strconv.Atoi(mux.Vars(r)[k])
	if err != nil {
		panic(fmt.Sprintf("mux var can't convert to int: %v", err))
	}
	return i
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		appErr, ok := errs.Cause(err).(*errs.Error)
		code := http.StatusInternalServerError
		if ok {
			code = appErr.HTTPCode
		}
		gores.String(w, code, fmt.Sprintf("%+v", err))
	}
}
