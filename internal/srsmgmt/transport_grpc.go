package srsmgmt

import (
	"context"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	grpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/gofrs/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"srsmgmt/config"
	pb "srsmgmt/pb"
)

type grpcServer struct {
	getStream    grpc.Handler
	createStream grpc.Handler
	deleteStream grpc.Handler
	pb.UnimplementedSrsMgmtServer
}

func MakeGRPCHandler(s Service, logger log.Logger) pb.SrsMgmtServer {
	e := MakeServerEndpoints(s)
	options := []grpc.ServerOption{
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		grpc.ServerBefore(jwt.GRPCToContext()),
	}

	cfgApikey := config.GetConfig().ApiKey

	return &grpcServer{
		getStream: grpc.NewServer(
			AuthMiddlewareGRPC(cfgApikey)(e.GetStreamEndpoint),
			decodeGRPCGetStreamRequest,
			encodeGRPCGetStreamResponse,
			options...,
		),
		createStream: grpc.NewServer(
			AuthMiddlewareGRPC(cfgApikey)(e.CreateStreamEndpoint),
			decodeGRPCCreateStreamRequest,
			encodeGRPCCreateStreamResponse,
			options...,
		),
		deleteStream: grpc.NewServer(
			AuthMiddlewareGRPC(cfgApikey)(e.DeleteStreamEndpoint),
			decodeGRPCDeleteStreamRequest,
			encodeGRPCDeleteStreamResponse,
			options...,
		),
	}
}

func (s *grpcServer) GetStream(ctx context.Context, req *pb.GetStreamRequest) (*pb.GetStreamReply, error) {
	_, rep, err := s.getStream.ServeGRPC(ctx, req)
	if err != nil {
		switch err {
		case ErrUnauthorized:
			return nil, status.Errorf(codes.Unauthenticated, ErrUnauthorized.Error())
		case ErrNotFound:
			return nil, status.Errorf(codes.NotFound, ErrNotFound.Error())
		default:
			return nil, status.Errorf(codes.InvalidArgument, ErrBadRequest.Error())
		}
	}
	return rep.(*pb.GetStreamReply), nil
}

func (s *grpcServer) CreateStream(ctx context.Context, req *pb.CreateStreamRequest) (*pb.CreateStreamReply, error) {
	_, rep, err := s.createStream.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateStreamReply), nil
}

func (s *grpcServer) DeleteStream(ctx context.Context, req *pb.DeleteStreamRequest) (*pb.DeleteStreamReply, error) {
	_, rep, err := s.deleteStream.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DeleteStreamReply), nil
}

func decodeGRPCGetStreamRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetStreamRequest)

	streamId, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, ErrBadRequest
	}

	endpointReq := getStreamRequest{}
	endpointReq.Stream.StreamID = streamId
	return endpointReq, nil
}

func encodeGRPCGetStreamResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getStreamResponse)

	if resp.Stream == nil {
		return nil, ErrNotFound
	}

	n := time.Time{}
	if resp.Stream.StartedAt == nil {
		resp.Stream.StartedAt = &n
	}
	if resp.Stream.StopedAt == nil {
		resp.Stream.StopedAt = &n
	}
	stream := &pb.Stream{
		Id:        resp.Stream.StreamID.String(),
		App:       resp.Stream.App,
		Password:  resp.Stream.Password,
		Status:    int32(resp.Stream.Status),
		Hls:       resp.Stream.HLS,
		CreatedAt: timestamppb.New(resp.Stream.CreatedAt),
		UpdatedAt: timestamppb.New(resp.Stream.UpdatedAt),
		ClientId:  resp.Stream.ClientId,
		StartedAt: timestamppb.New(*resp.Stream.StartedAt),
		StopedAt:  timestamppb.New(*resp.Stream.StopedAt),
	}

	return &pb.GetStreamReply{Stream: stream}, nil
}

func decodeGRPCCreateStreamRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.CreateStreamRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	streamId, err := uuid.FromString(req.Stream.Id)
	if err != nil {
		return nil, ErrBadRequest
	}

	stream := Stream{
		StreamID: streamId,
		App:      req.Stream.App,
		Password: req.Stream.Password,
	}

	return createStreamRequest{Stream: stream}, nil
}

func encodeGRPCCreateStreamResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	resp, ok := grpcReply.(createStreamResponse)
	if !ok {
		return nil, ErrBadRequest
	}

	n := time.Time{}
	if resp.Stream.StartedAt == nil {
		resp.Stream.StartedAt = &n
	}
	if resp.Stream.StopedAt == nil {
		resp.Stream.StopedAt = &n
	}
	stream := &pb.Stream{
		Id:        resp.Stream.StreamID.String(),
		App:       resp.Stream.App,
		Password:  resp.Stream.Password,
		Status:    int32(resp.Stream.Status),
		Hls:       resp.Stream.HLS,
		CreatedAt: timestamppb.New(resp.Stream.CreatedAt),
		UpdatedAt: timestamppb.New(resp.Stream.UpdatedAt),
		ClientId:  resp.Stream.ClientId,
		StartedAt: timestamppb.New(*resp.Stream.StartedAt),
		StopedAt:  timestamppb.New(*resp.Stream.StopedAt),
	}

	return &pb.CreateStreamReply{Stream: stream}, nil
}

func decodeGRPCDeleteStreamRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteStreamRequest)

	streamId, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, ErrBadRequest
	}

	endpointReq := deleteStreamRequest{}
	endpointReq.ID = streamId
	return endpointReq, nil
}

func encodeGRPCDeleteStreamResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(deleteStreamResponse)

	// if resp.ID == nil {
	// 	return &pb.DeleteStreamReply{Id: resp.ID.String(), Err: "1"}, nil
	// }
	return &pb.DeleteStreamReply{Id: resp.ID.String()}, nil
}
