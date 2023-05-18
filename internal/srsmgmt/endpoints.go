package srsmgmt

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofrs/uuid"
)

type Endpoints struct {
	GetStreamEndpoint       endpoint.Endpoint
	CreateStreamEndpoint    endpoint.Endpoint
	DeleteStreamEndpoint    endpoint.Endpoint
	StartStreamEndpoint     endpoint.Endpoint
	StopStreamEndpoint      endpoint.Endpoint
	MonStreamEndpoint       endpoint.Endpoint
	UpdateSRSStreamEndpoint endpoint.Endpoint
	UpdateHlsSRSEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetStreamEndpoint:       MakeGetStreamEndpoint(s),
		CreateStreamEndpoint:    MakeCreateStreamEndpoint(s),
		DeleteStreamEndpoint:    MakeDeleteStreamEndpoint(s),
		StartStreamEndpoint:     MakeStartStreamEndpoint(s),
		StopStreamEndpoint:      MakeStopStreamEndpoint(s),
		MonStreamEndpoint:       MakeMonStreamEndpoint(s),
		UpdateSRSStreamEndpoint: MakeUpdateSRSStreamEndpoint(s),
		UpdateHlsSRSEndpoint:    MakeUpdateHlsSRSEndpoint(s),
	}
}

func MakeGetStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getStreamRequest)
		stream, e := s.GetStream(ctx, req.Stream)

		return getStreamResponse{Stream: stream}, e
	}
}

func MakeDeleteStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteStreamRequest)
		streamId, e := s.DeleteStream(ctx, req.ID)
		return deleteStreamResponse{ID: &streamId}, e
	}
}

func MakeStartStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(startStreamRequest)
		stream, e := s.StartStream(ctx, req.ID)
		return startStreamResponse{Stream: stream}, e
	}
}

func MakeStopStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(stopStreamRequest)
		stream, e := s.StopStream(ctx, req.ID)
		return stopStreamResponse{Stream: stream}, e
	}
}

func MakeMonStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		stream, e := s.MonStream(ctx)
		return stream, e
	}
}

func MakeCreateStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createStreamRequest)
		stream, e := s.CreateStream(ctx, req.Stream)

		return createStreamResponse{Stream: stream}, e
	}
}

func MakeUpdateSRSStreamEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateSRSStreamRequest)
		code, e := s.UpdateSRSStream(ctx, req.Stream)

		return updateSRSStreamResponse{Code: code}, e
	}
}

func MakeUpdateHlsSRSEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateHlsSRSRequest)
		code, e := s.UpdateHlsSRS(ctx, req.Stream)

		return updateHlsSRSResponse{Code: code}, e
	}
}

type getStreamRequest struct {
	Stream Stream `json:"stream,omitempty"`
}

type getStreamResponse struct {
	Stream *Stream `json:"stream,omitempty"`
}

type createStreamRequest struct {
	Stream Stream `json:"stream,omitempty"`
}

type createStreamResponse struct {
	Stream *Stream `json:"stream,omitempty"`
}

type updateSRSStreamRequest struct {
	Stream SRSStream `json:"stream,omitempty"`
}

type updateSRSStreamResponse struct {
	Code int `json:"code"`
}

type updateHlsSRSRequest struct {
	Stream SRSStream `json:"stream,omitempty"`
}

type updateHlsSRSResponse struct {
	Code int `json:"code"`
}

type deleteStreamRequest struct {
	ID uuid.UUID `json:"streamId"`
}

type deleteStreamResponse struct {
	ID *uuid.UUID `json:"stream,omitempty"`
}

type stopStreamRequest struct {
	ID uuid.UUID `json:"streamId"`
}

type stopStreamResponse struct {
	Stream *Stream `json:"stream,omitempty"`
}

type monStreamResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	App   string `json:"app"`
	Video struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"video"`
	Audio struct {
		Channel int `json:"channel"`
	} `json:"audio"`
	Kbps struct {
		Recv_30s int `json:"recv_30s"`
	}
}

type startStreamRequest struct {
	ID uuid.UUID `json:"streamId"`
}

type startStreamResponse struct {
	Stream *Stream `json:"stream,omitempty"`
}
