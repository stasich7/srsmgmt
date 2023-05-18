package srsclient

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Middleware func(SrsClient) SrsClient

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next SrsClient) SrsClient {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   SrsClient
	logger log.Logger
}

func (mw loggingMiddleware) GetSRSStream() (p *[]SRSStream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Client GetSRSStreams", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetSRSStream()
}

func (mw loggingMiddleware) KickSRSStream(s string) (err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Client KickSRSStream", "id", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.KickSRSStream(s)
}

func (mw loggingMiddleware) ConfigReload() (err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Client ConfigReload", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ConfigReload()
}
