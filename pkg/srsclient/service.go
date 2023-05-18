package srsclient

import (
	"errors"

	"github.com/go-kit/log"
)

var (
	ErrInternalError = errors.New("INTERNAL_ERROR")
)

type SrsClient interface {
	GetSRSStream() (*[]SRSStream, error)
	KickSRSStream(string) error
	ConfigReload() error
}

type srsClientService struct {
	logger   log.Logger
	clntSrvc SrsClientEndpoint
}

type SRSPublish struct {
	Active bool   `json:"active"`
	Cid    string `json:"cid"`
}

type SRSStream struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	App     string      `json:"app"`
	Publish *SRSPublish `json:"publish"`
	Video   struct {
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

func New(serverURL string, logger log.Logger) SrsClient {
	return NewBasicService(serverURL, logger)
}

func NewBasicService(serverURL string, logger log.Logger) *srsClientService {
	client := NewEndpoints(serverURL)
	return &srsClientService{logger, client}
}

func (sc *srsClientService) ConfigReload() error {
	err := sc.clntSrvc.ConfigReload()
	if err != nil {
		return nil
	}
	return nil
}

func (sc *srsClientService) GetSRSStream() (*[]SRSStream, error) {
	resp, err := sc.clntSrvc.GetSRSStreams()
	if err != nil {
		return nil, err
	}

	liveStreams := []SRSStream{}
	for i, v := range *resp.Streams {
		if v.App == "live" && v.Publish.Active && indexOf(v, liveStreams) == -1 {
			(*resp.Streams)[i].Kbps.Recv_30s = (*resp.Streams)[i].Kbps.Recv_30s << 4
			liveStreams = append(liveStreams, v)
		}
	}
	return &liveStreams, nil
}

func (sc *srsClientService) KickSRSStream(cid string) error {
	if cid == "" {
		return nil
	}

	err := sc.clntSrvc.KickSRSStream(cid)
	if err != nil {
		return nil
	}
	return nil
}

func indexOf(item SRSStream, stor []SRSStream) int {
	for i, v := range stor {
		if v.Id == item.Id {
			return i
		}
	}
	return -1
}
