package srsmgmt

import (
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"srsmgmt/config"
	"srsmgmt/pkg/playlist"
	"srsmgmt/pkg/srsclient"
	"srsmgmt/pkg/srsconfig"
	"strings"
	"sync"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"
)

const (
	StreamStatusWaitPublish   = 1
	StreamStatusPublish       = 2
	StreamStatusPause         = 3
	StreamStatusStopPublish   = 4
	StreamStatusError         = 5
	StreamStatusStartRequired = 6
	SRSok                     = 0
	SRSfail                   = 1
	OutputPlaylistPrefix      = "master-"
)

var (
	ErrNotFound      = errors.New("STREAM_NOT_FOUND")
	ErrAlreadyExists = errors.New("ALREADY_EXISTS")
	ErrNoStreaming   = errors.New("STREAMING_NOT_STARTED")
	ErrBadRequest    = errors.New("BAD_REQUEST")
	ErrBadStatus     = errors.New("BAD_STATUS")
	ErrInternalError = errors.New("INTERNAL_ERROR")
	ErrUnauthorized  = errors.New("UNAUTHORIZED")
)

type Service interface {
	GetStream(context.Context, Stream) (*Stream, error)
	CreateStream(context.Context, Stream) (*Stream, error)
	DeleteStream(context.Context, uuid.UUID) (uuid.UUID, error)
	StopStream(context.Context, uuid.UUID) (*Stream, error)
	StartStream(context.Context, uuid.UUID) (*Stream, error)
	MonStream(context.Context) (*[]monStream, error)
	UpdateSRSStream(context.Context, SRSStream) (int, error)
	UpdateHlsSRS(context.Context, SRSStream) (int, error)
}

type Repository interface {
	GetMock() sqlmock.Sqlmock
	GetStream(uuid.UUID) (*Stream, error)
	GetRTCStreams() (*[]Stream, error)
	CreateStream(Stream) (*Stream, error)
	DeleteStream(uuid.UUID) (uuid.UUID, error)
	UpdateStream(Stream) (*Stream, error)
}

// swagger:model Stream
type Stream struct {
	StreamID  uuid.UUID  `json:"id"`
	App       string     `json:"app"`
	Password  string     `json:"password"`
	Status    int        `json:"status"`
	HLS       string     `json:"hls"`
	RTMPpush  string     `json:"rtmpPush"`
	SRTpush   string     `json:"srtPush"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ClientId  string     `json:"clientId"`
	StartedAt *time.Time `json:"startedAt"`
	StopedAt  *time.Time `json:"stopedAt"`
	RTC       bool       `json:"rtc"`
}

type SRSStream struct {
	Action   string  `json:"action"`
	ClientID string  `json:"client_id"`
	IP       string  `json:"ip"`
	Vhost    string  `json:"vhost"`
	App      string  `json:"app"`
	StreamID string  `json:"stream"`
	Param    string  `json:"param"`
	Duration float32 `json:"duration"`
	Cwd      string  `json:"cwd"`
	File     string  `json:"file"`
	URL      string  `json:"url"`
	M3U8     string  `json:"m3u8"`
	M3U8_URL string  `json:"m3u8_url"`
	Seq_no   int     `json:"seq_no"`
}

type cachedStreams struct {
	streams []monStream
	time.Time
	sync.Mutex
}

type srsMgmtService struct {
	cfg           config.Config
	repo          Repository
	logger        log.Logger
	playlist      *playlist.Playlist
	srsClient     srsclient.SrsClient
	srsConfig     *srsconfig.SRSConfig
	cachedStreams cachedStreams
}

type monStream struct {
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

func NewSrsMgmtService(repo Repository, logger log.Logger, plist *playlist.Playlist, client srsclient.SrsClient, srsconfigSvc *srsconfig.SRSConfig) Service {
	rtcStreamsRepo, _ := repo.GetRTCStreams()
	rtcStreams := []srsconfig.SRSstream{}
	for _, v := range *rtcStreamsRepo {
		rtcStreams = append(rtcStreams, srsconfig.SRSstream{
			ID:       v.StreamID.String(),
			Password: v.Password,
		})
	}
	srsconfigSvc.Init(rtcStreams)

	s := srsMgmtService{
		cfg:       *config.GetConfig(),
		repo:      repo,
		logger:    logger,
		playlist:  plist,
		srsClient: client,
		srsConfig: srsconfigSvc,
	}
	return &s
}

func (s *srsMgmtService) GetStream(ctx context.Context, st Stream) (*Stream, error) {
	stream, err := s.repo.GetStream(st.StreamID)
	if err != nil {
		return &Stream{}, ErrNotFound
	}
	s.addSRSUrls(stream)

	return stream, err
}

func (s *srsMgmtService) CreateStream(ctx context.Context, newStream Stream) (*Stream, error) {
	// временно все стримы становятся webrtc
	newStream.RTC = true

	stream, err := s.repo.GetStream(newStream.StreamID)
	if err != nil || stream == nil {
		// если стрима нет, то создаем
		stream, err = s.repo.CreateStream(newStream)
		if err != nil {
			return &Stream{}, ErrInternalError
		}
	} else {
		// если стрим есть, то "активируем его" , переведя в статус StreamStatusWaitPublish
		stream.Status = StreamStatusWaitPublish
		stream.StopedAt = nil
		stream.StartedAt = nil
		stream.RTC = newStream.RTC
		stream.Password = newStream.Password
		stream, err = s.repo.UpdateStream(*stream)
		if err != nil {
			return nil, ErrInternalError
		}
	}
	s.addSRSUrls(stream)

	if stream.RTC {
		s.srsConfig.AddRTC(stream.StreamID.String(), stream.Password)
	} else {
		s.srsConfig.RemoveRTC(stream.StreamID.String())
	}

	return stream, err
}

func (s *srsMgmtService) DeleteStream(ctx context.Context, streamID uuid.UUID) (uuid.UUID, error) {
	stream, err := s.repo.GetStream(streamID)
	if err != nil {
		return streamID, ErrNotFound
	}

	s.logger.Log("Kicking client ", stream.ClientId)
	if err := s.srsClient.KickSRSStream(stream.ClientId); err != nil {
		return uuid.UUID{}, ErrInternalError
	}

	if err := s.playlist.Delete(path.Join(s.cfg.LiveTSPath, stream.StreamID.String())); err != nil {
		return uuid.UUID{}, ErrInternalError
	}

	s.srsConfig.RemoveRTC(stream.StreamID.String())

	_, err = s.repo.DeleteStream(stream.StreamID)
	if err != nil {
		return uuid.UUID{}, ErrInternalError
	}

	return streamID, err
}

func (s *srsMgmtService) StartStream(ctx context.Context, streamID uuid.UUID) (*Stream, error) {
	stream, err := s.repo.GetStream(streamID)
	if err != nil {
		return nil, ErrNotFound
	}

	if stream.Status == StreamStatusStopPublish {
		return nil, ErrBadStatus
	}

	n := time.Now()
	if stream.Status == StreamStatusWaitPublish || stream.Status == StreamStatusStartRequired {
		stream.Status = StreamStatusStartRequired
		stream.StartedAt = &n
		stream, err = s.repo.UpdateStream(*stream)
		if err != nil {
			return nil, ErrInternalError
		}
		s.addSRSUrls(stream)
		return stream, nil
	}
	stream.StartedAt = &n

	level.Debug(s.logger).Log("before playlist.StartStream:Refresh")
	if err := s.playlist.Refresh(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), playlist.AllPlaylists, stream.StartedAt, OutputPlaylistPrefix); err != nil {
		level.Debug(s.logger).Log("playlist.StartStream:Refresh", err)
		return nil, err
	}

	level.Debug(s.logger).Log("before playlist.StartStream:Create", err)
	if err := s.playlist.Create(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), OutputPlaylistPrefix, stream.StartedAt); err != nil {
		level.Debug(s.logger).Log("playlist.StartStream:Create", err)
	}
	stream, err = s.repo.UpdateStream(*stream)
	if err != nil {
		return nil, ErrInternalError
	}
	s.addSRSUrls(stream)

	return stream, err
}

func (s *srsMgmtService) StopStream(ctx context.Context, streamID uuid.UUID) (*Stream, error) {
	stream, err := s.repo.GetStream(streamID)
	if err != nil {
		return nil, ErrNotFound
	}

	if stream.Status == StreamStatusStopPublish {
		return nil, ErrBadStatus
	}

	level.Info(s.logger).Log("Kicking client ", stream.ClientId)
	if err := s.srsClient.KickSRSStream(stream.ClientId); err != nil {
		return nil, err
	}

	if err := s.playlist.Stop(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), OutputPlaylistPrefix, stream.StartedAt); err != nil {
		return nil, ErrBadStatus
	}

	n := time.Now()
	stream.Status = StreamStatusStopPublish
	stream.ClientId = ""
	stream.StopedAt = &n
	stream.RTC = false
	stream, err = s.repo.UpdateStream(*stream)
	if err != nil {
		return nil, ErrInternalError
	}
	s.addSRSUrls(stream)

	s.srsConfig.RemoveRTC(stream.StreamID.String())

	return stream, err
}

func (s *srsMgmtService) MonStream(ctx context.Context) (*[]monStream, error) {
	s.cachedStreams.Lock()
	defer s.cachedStreams.Unlock()

	if s.cachedStreams.Add(time.Duration(s.cfg.CacheTTL) * time.Second).Before(time.Now()) {
		streamStat, err := s.srsClient.GetSRSStream()
		if err != nil || streamStat == nil {
			return nil, err
		}

		s.cachedStreams.streams = []monStream{}
		for _, v := range *streamStat {
			s.cachedStreams.streams = append(s.cachedStreams.streams, monStream{
				Id:    v.Id,
				Name:  v.Name,
				App:   v.App,
				Video: v.Video,
				Audio: v.Audio,
				Kbps:  v.Kbps,
			})
		}
		s.cachedStreams.Time = time.Now()
	}

	return &s.cachedStreams.streams, nil
}

func (s *srsMgmtService) UpdateSRSStream(ctx context.Context, st SRSStream) (int, error) {
	if strings.Contains(st.Param, "upstream=rtc") {
		return SRSok, nil
	}
	if st.App != "live" {
		return SRSfail, ErrBadRequest
	}
	streamUuid, err := uuid.FromString(st.StreamID)
	if err != nil {
		return SRSfail, ErrNotFound
	}
	re := regexp.MustCompile("password=([^&$]+)")
	passParams := re.FindStringSubmatch(st.Param)
	if len(passParams) == 0 {
		return SRSfail, ErrBadRequest
	}
	stream, err := s.repo.GetStream(streamUuid)
	if err != nil {
		return SRSfail, ErrNotFound
	}
	if stream.Password != passParams[1] {
		return SRSfail, ErrUnauthorized
	}

	switch st.Action {
	case "on_publish":
		switch stream.Status {
		case StreamStatusStartRequired:
			go func(stream *Stream) {
				if err := s.playlist.Create(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), OutputPlaylistPrefix, stream.StartedAt); err != nil {
					level.Debug(s.logger).Log("playlist.CreateDVR", err)
					stream.Status = StreamStatusError
					s.repo.UpdateStream(*stream)

					level.Debug(s.logger).Log("Kicking client ", stream.ClientId)
					s.srsClient.KickSRSStream(stream.ClientId)
				}
			}(stream)
			fallthrough

		case StreamStatusWaitPublish, StreamStatusError:
			stream.Status = StreamStatusPublish
			stream.ClientId = st.ClientID
			s.repo.UpdateStream(*stream)
			go func(stream *Stream) {
				if err := s.playlist.Create(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), OutputPlaylistPrefix, nil); err != nil {
					level.Debug(s.logger).Log("playlist.Create", err)
					stream.Status = StreamStatusError
					s.repo.UpdateStream(*stream)

					level.Debug(s.logger).Log("Kicking client ", stream.ClientId)
					s.srsClient.KickSRSStream(stream.ClientId)
				}
			}(stream)

			return SRSok, nil

		case StreamStatusPause, StreamStatusPublish:
			stream.Status = StreamStatusPublish
			stream.ClientId = st.ClientID

		default:
			return SRSfail, ErrBadRequest
		}
	case "on_unpublish":
		switch stream.Status {
		case StreamStatusPublish:
			stream.Status = StreamStatusPause
			stream.ClientId = ""
		default:
			return SRSok, nil
		}
	default:
		return SRSfail, ErrBadRequest
	}
	_, err = s.repo.UpdateStream(*stream)
	if err != nil {
		return SRSfail, err
	}
	return SRSok, nil
}

func (s *srsMgmtService) UpdateHlsSRS(ctx context.Context, st SRSStream) (int, error) {
	if strings.Contains(st.Param, "upstream=rtc") {
		return SRSok, nil
	}

	streamUuid, err := uuid.FromString(st.StreamID)
	if err != nil {
		return SRSfail, ErrNotFound
	}

	stream, err := s.repo.GetStream(streamUuid)
	if err != nil {
		return SRSfail, ErrNotFound
	}

	plName := path.Base(st.M3U8)
	if stream.Status != StreamStatusStopPublish {
		if err := s.playlist.Refresh(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), plName, stream.StartedAt, OutputPlaylistPrefix); err != nil {
			level.Debug(s.logger).Log("Refresh Playlist", fmt.Sprintf("[%s] %s", plName, stream.StartedAt))
			return 0, err
		}
	}

	if stream.Status == StreamStatusStartRequired {
		go func(stream *Stream) {
			stream.Status = StreamStatusPublish
			s.repo.UpdateStream(*stream)

			if err := s.playlist.Create(path.Join(s.cfg.LiveTSPath, stream.StreamID.String()), OutputPlaylistPrefix, stream.StartedAt); err != nil {
				level.Debug(s.logger).Log("playlist.CreateDVR", err)
				stream.Status = StreamStatusError
				s.repo.UpdateStream(*stream)
			}
		}(stream)
	}

	return 0, nil
}

func (s *srsMgmtService) addSRSUrls(stream *Stream) {
	streamTS := ""
	if stream.StartedAt != nil && !stream.StartedAt.IsZero() {
		streamTS = fmt.Sprintf("%d-", stream.StartedAt.Unix())
	}
	stream.HLS = fmt.Sprintf("%s/%s/%s/%s%s%s", s.cfg.HLSAddr, stream.App, stream.StreamID.String(), OutputPlaylistPrefix, streamTS, "index.m3u8")
	stream.RTMPpush = fmt.Sprintf("%s/%s/%s?password=%s", s.cfg.RTMPAddr, stream.App, stream.StreamID.String(), stream.Password)
	stream.SRTpush = fmt.Sprintf("%s?streamid=#!::r=%s/%s,m=publish,password=%s", s.cfg.SRTAddr, stream.App, stream.StreamID.String(), stream.Password)
}
