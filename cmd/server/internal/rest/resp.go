package rest

import (
	"boiler/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrResponse struct {
	Error struct {
		Codes []string `json:"codes"`
		Msg   string   `json:"msg"`
	} `json:"error"`
}

type DefaultResp struct{}

// Fail writes the JSON error message
// if ?debug is set it will response with the original errors message
func (d DefaultResp) Fail(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	resp := new(ErrResponse)

	resp.Error.Msg = err.Error()
	if errors.Is(err, errors.ErrBadRequest) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Error().Err(err).Str("file", errors.Caller()).Send()
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error.Codes = append(resp.Error.Codes, "internal_server_error")

		if len(r.URL.Query()["debug"]) == 0 {
			resp.Error.Msg = http.StatusText(http.StatusInternalServerError)
		}
	}

	var ce *errors.CodeErr
	errs := err
	for {
		if errors.As(errs, &ce) {
			resp.Error.Codes = append(resp.Error.Codes, ce.Code)
			errs = errors.Unwrap(ce)
			continue
		}

		break
	}

	d.JSON(w, r, resp)
}

// FailF same as Fail, but with error format
func (d DefaultResp) Failf(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) {
	d.Fail(w, r, fmt.Errorf(format, a...))
}

// JSON writes the content of the param data as JSON.
// if ?pretty is present, it will pretty print the response.
func (d DefaultResp) JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	if len(r.URL.Query()["pretty"]) != 0 {
		e.SetIndent(" ", " ")
	}
	if err := e.Encode(data); err != nil {
		d.Fail(w, r, fmt.Errorf("could not write json response; %w", err))
	}
}
