package srsmgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"srsmgmt/config"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	g := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	}

	cfgApikey := config.GetConfig().ApiKey

	r := g.PathPrefix("/api/v1/").Subrouter()

	r.Methods("GET").Path("/stream/{id}").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.GetStreamEndpoint),
		decodeGetStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/stream/{id}").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.CreateStreamEndpoint),
		decodeCreateStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/stream/{id}").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.DeleteStreamEndpoint),
		decodeDeleteStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/stream/{id}/start").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.StartStreamEndpoint),
		decodeStartStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/stream/{id}/stop").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.StopStreamEndpoint),
		decodeStopStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/streams/monitor").Handler(httptransport.NewServer(
		AuthMiddlewareHTTP(cfgApikey)(e.MonStreamEndpoint),
		decodeMonStreamRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/webhook/stream/live").Handler(httptransport.NewServer(
		e.UpdateSRSStreamEndpoint,
		decodeUpdateSRSStreamRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/webhook/stream/hls").Handler(httptransport.NewServer(
		e.UpdateHlsSRSEndpoint,
		decodeUpdateHlsSRSRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRequest
	}
	uuid, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	var req getStreamRequest
	req.Stream.StreamID = uuid
	return req, nil
}

func decodeDeleteStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRequest
	}
	streamId, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	var req deleteStreamRequest
	req.ID = streamId
	return req, nil
}
func decodeStartStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRequest
	}
	streamId, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	var req startStreamRequest
	req.ID = streamId
	return req, nil
}

func decodeStopStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRequest
	}
	streamId, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	var req stopStreamRequest
	req.ID = streamId
	return req, nil
}

func decodeMonStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return
}

func decodeUpdateSRSStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req updateSRSStreamRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Stream); e != nil {
		return nil, ErrBadRequest
	}
	return req, nil
}

func decodeUpdateHlsSRSRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req updateHlsSRSRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Stream); e != nil {
		return nil, ErrBadRequest
	}
	return req, nil
}

func decodeCreateStreamRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRequest
	}
	streamId, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	var req createStreamRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Stream); e != nil {
		return nil, e
	}
	req.Stream.StreamID = streamId
	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusUnprocessableEntity
	case ErrNotFound, ErrBadRequest, ErrBadStatus:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
