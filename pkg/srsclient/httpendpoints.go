package srsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
)

const STREAMS_COUNT = 100 // лимит стримов для мониторинга
type SrsClientEndpoint interface {
	GetSRSStreams() (*GetSRSStreamsResponse, error)
	KickSRSStream(string) error
	ConfigReload() error
}
type SrsClientSet struct {
	GetSRSStreamEndpoint  endpoint.Endpoint
	KickSRSStreamEndpoint endpoint.Endpoint
	ConfigReloadEndPoint  endpoint.Endpoint
}

func NewEndpoints(serverAddr string) SrsClientSet {
	setSRSStreamTransport := HttpRequest{
		Addr:    serverAddr,
		Method:  "GET",
		Encoder: encodeGetSRSStreamRequest,
		Decoder: decodeGetSRSStreamsResponse,
	}

	setSRSKickTransport := HttpRequest{
		Addr:    serverAddr,
		Method:  "DELETE",
		Encoder: encodeKickSRSStreamRequest,
		Decoder: decodeGetSRSStreamsResponse,
	}

	setConfigReloadTransport := HttpRequest{
		Addr:    serverAddr,
		Method:  "GET",
		Encoder: encodeConfigReloadRequest,
		Decoder: decodeConfigReloadResponse,
	}

	return SrsClientSet{
		GetSRSStreamEndpoint:  setSRSStreamTransport.MakeRequest("/api/v1/streams"),
		KickSRSStreamEndpoint: setSRSKickTransport.MakeRequest("/api/v1/clients/"),
		ConfigReloadEndPoint:  setConfigReloadTransport.MakeRequest("/api/v1/raw"),
	}
}

type GetSRSStreamsRequest struct {
}

type GetSRSStreamsResponse struct {
	Code    int          `json:"code"`
	Streams *[]SRSStream `json:"streams"`
	Error   string       `json:"error"`
}

type KickSRSStreamRequest struct {
	Cid string
}

type GetConfigReloadRequest struct {
}

type GetConfigReloadResponse struct {
}

func (s SrsClientSet) GetSRSStreams() (*GetSRSStreamsResponse, error) {
	requestData := GetSRSStreamsRequest{}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := s.GetSRSStreamEndpoint(ctx, requestData)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New(response.(GetSRSStreamsResponse).Error)
	}

	resp, ok := response.(GetSRSStreamsResponse)
	if !ok {
		return nil, errors.New(response.(GetSRSStreamsResponse).Error)
	}

	return &resp, nil
}

func (s SrsClientSet) ConfigReload() error {
	requestData := GetConfigReloadRequest{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := s.ConfigReloadEndPoint(ctx, requestData)
	fmt.Printf("ConfigReload err:%v response:%v\n", err, response)
	if err != nil {
		return err
	}

	if response == nil {
		return errors.New(response.(GetSRSStreamsResponse).Error)
	}
	return nil
}

func (s SrsClientSet) KickSRSStream(cid string) error {
	requestData := KickSRSStreamRequest{Cid: cid}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := s.KickSRSStreamEndpoint(ctx, requestData)
	if err != nil {
		return err
	}

	if response == nil {
		return errors.New(response.(GetSRSStreamsResponse).Error)
	}
	return nil
}

func decodeGetSRSStreamsResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}
	var response GetSRSStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func encodeGetSRSStreamRequest(ctx context.Context, req *http.Request, request interface{}) error {
	q := req.URL.Query()
	q.Add("count", fmt.Sprintf("%d", STREAMS_COUNT))
	req.URL.RawQuery = q.Encode()
	return nil
}

func encodeConfigReloadRequest(ctx context.Context, req *http.Request, request interface{}) error {
	q := req.URL.Query()
	q.Add("rpc", "reload")
	req.URL.RawQuery = q.Encode()
	return nil
}

func decodeConfigReloadResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	return resp, nil
}

func encodeKickSRSStreamRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(KickSRSStreamRequest)
	Cid := url.QueryEscape(r.Cid)
	req.URL.Path = fmt.Sprintf("%s%s", req.URL.Path, Cid)
	return nil
}
