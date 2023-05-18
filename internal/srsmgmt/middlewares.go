package srsmgmt

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/metadata"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetStream(ctx context.Context, s Stream) (p *Stream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "GetStream", "id", s.StreamID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetStream(ctx, s)
}

func (mw loggingMiddleware) DeleteStream(ctx context.Context, s uuid.UUID) (id uuid.UUID, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "DeleteStream", "id", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteStream(ctx, s)
}

func (mw loggingMiddleware) StartStream(ctx context.Context, s uuid.UUID) (p *Stream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "StartStream", "id", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.StartStream(ctx, s)
}

func (mw loggingMiddleware) StopStream(ctx context.Context, s uuid.UUID) (p *Stream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "StopStream", "id", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.StopStream(ctx, s)
}

func (mw loggingMiddleware) MonStream(ctx context.Context) (p *[]monStream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "MonStream", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.MonStream(ctx)
}

func (mw loggingMiddleware) CreateStream(ctx context.Context, s Stream) (p *Stream, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "CreateStream", "data", fmt.Sprintf("%+v", s), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateStream(ctx, s)
}

func (mw loggingMiddleware) UpdateSRSStream(ctx context.Context, s SRSStream) (code int, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "UpdateSRSStream", "data", fmt.Sprintf("%+v", s), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateSRSStream(ctx, s)
}

func (mw loggingMiddleware) UpdateHlsSRS(ctx context.Context, s SRSStream) (code int, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "UpdateHlsSRS", "data", fmt.Sprintf("%+v", s), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateHlsSRS(ctx, s)
}

func AuthMiddlewareHTTP(requiredApiKey string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			auth, ok := ctx.Value(httptransport.ContextKeyRequestAuthorization).(string)
			if !ok || auth != requiredApiKey {
				return nil, ErrUnauthorized
			}
			return next(ctx, request)
		}
	}
}

func AuthMiddlewareGRPC(requiredApiKey string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			auth := md.Get("authorization")
			if !ok || auth == nil || auth[0] != requiredApiKey {
				return nil, ErrUnauthorized
			}
			return next(ctx, request)
		}
	}
}
