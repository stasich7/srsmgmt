package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"srsmgmt/config"
	"srsmgmt/internal/srsmgmt"
	"srsmgmt/internal/srsmgmtrepo"
	pb "srsmgmt/pb"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	glogger "gorm.io/gorm/logger"
)

func TestHTTP(t *testing.T) {
	logger := log.NewNopLogger()
	repo := srsmgmtrepo.New(logger, glogger.Silent, config.GetConfig())

	svc := srsmgmt.NewSrsMgmtService(repo, logger, nil, nil)
	httpHandler := srsmgmt.MakeHTTPHandler(svc, logger)
	srv := httptest.NewServer(httpHandler)

	defer srv.Close()

	for _, testcase := range []struct {
		method, url, body, want string
		statuscode              int
	}{
		{"DELETE", srv.URL + "/api/v1/stream/00000000-1111-0000-0000-000000000000", ``, ``, 0},
		{"POST", srv.URL + "/api/v1/stream/00000000-1111-0000-0000-000000000000", `{"app": "live","password": "123"}`, `00000000-1111-0000-0000-000000000000`, http.StatusOK},
		{"GET", srv.URL + "/api/v1/stream/00000000-1111-0000-0000-000000000000", ``, `00000000-1111-0000-0000-000000000000`, http.StatusOK},
		{"POST", srv.URL + "/api/v1/webhook/stream/live", `{"action":"on_publish","client_id":"","ip":"","vhost":"","app":"live","stream": "00000000-1111-0000-0000-000000000000","fail": false,"param": "?password=viod"}`, `UNAUTHORIZED`, http.StatusUnauthorized},
		{"DELETE", srv.URL + "/api/v1/stream/00000000-1111-0000-0000-000000000000", ``, `00000000-1111-0000-0000-000000000000`, http.StatusOK},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		req.Header.Add("Authorization", config.GetConfig().ApiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			break
		}

		if status := resp.StatusCode; testcase.statuscode > 0 && status != testcase.statuscode {
			t.Errorf("handler [%s] returned wrong status code: got %v want %v", testcase.url, status, http.StatusOK)
			break
		}

		respBuf, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			break
		}
		// t.Logf(string(respBuf))

		if want, have := testcase.want, string(respBuf); !strings.Contains(have, want) {
			t.Errorf("\n%s %s %s: \nwant contain %q, \nhave %q", testcase.method, testcase.url, testcase.body, want, have)
		}
	}
}

func TestGRPC(t *testing.T) {
	logger := log.NewNopLogger()
	repo := srsmgmtrepo.New(logger, glogger.Silent, config.GetConfig())
	svc := srsmgmt.NewSrsMgmtService(repo, logger, nil, nil)
	grpcListener, err := net.Listen("tcp", config.GetConfig().GRPCAddr)
	if err != nil {
		t.Fatalf("GRPC Connection failure: %v", err)
	}
	grpcHandler := srsmgmt.MakeGRPCHandler(svc, logger)
	s := grpc.NewServer()
	pb.RegisterSrsMgmtServer(s, grpcHandler)

	go func() {
		s.Serve(grpcListener)
	}()

	defer grpcListener.Close()

	var address = config.GetConfig().GRPCAddr
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("GRPC Connection failure: %v", err)
	}
	defer conn.Close()
	c := pb.NewSrsMgmtClient(conn)

	testcase := []struct {
		name, requestId, requestApp, requestPassword, want, ignoreErr string
	}{
		{
			name:      "DeleteStream",
			requestId: "00000000-2222-0000-0000-000000000000",
			want:      "00000000-2222-0000-0000-000000000000",
			ignoreErr: "STREAM_NOT_FOUND",
		},
		{
			name:            "CreateStream",
			requestId:       "00000000-2222-0000-0000-000000000000",
			requestApp:      "live",
			requestPassword: "123",
			want:            "00000000-2222-0000-0000-000000000000",
		},
		{
			name:      "GetStream",
			requestId: "00000000-2222-0000-0000-000000000000",
			want:      "00000000-2222-0000-0000-000000000000",
		},
		{
			name:      "DeleteStream",
			requestId: "00000000-2222-0000-0000-000000000000",
			want:      "00000000-2222-0000-0000-000000000000",
		},
	}

	for _, tt := range testcase {
		// t.Log("Running", tt.name)

		switch tt.name {
		case "CreateStream":
			stream := &pb.Stream{Id: tt.requestId, App: tt.requestApp, Password: tt.requestPassword}
			req := &pb.CreateStreamRequest{Stream: stream}
			header := metadata.New(map[string]string{"authorization": config.GetConfig().ApiKey})
			ctx := metadata.NewOutgoingContext(context.Background(), header)

			resp, err := c.CreateStream(ctx, req)
			if err != nil {
				errStatus, _ := status.FromError(err)
				t.Errorf("%s got unexpected error %q", tt.name, errStatus.Message())
				break
			}
			if resp.Stream.Id != tt.want {
				t.Errorf("%v, \nwant contain %q, \nhave %q", tt.name, tt.want, resp.Stream.Id)
			}

		case "GetStream":
			req := &pb.GetStreamRequest{Id: tt.requestId}
			header := metadata.New(map[string]string{"authorization": config.GetConfig().ApiKey})
			ctx := metadata.NewOutgoingContext(context.Background(), header)

			resp, err := c.GetStream(ctx, req)
			if err != nil {
				errStatus, _ := status.FromError(err)
				t.Errorf("%s got unexpected error %q", tt.name, errStatus.Message())
				break
			}
			if resp.Stream.Id != tt.want {
				t.Errorf("%v, \nwant contain %q, \nhave %q", tt.name, tt.want, resp.Stream.Id)
			}

		case "DeleteStream":
			req := &pb.DeleteStreamRequest{Id: tt.requestId}
			header := metadata.New(map[string]string{"authorization": config.GetConfig().ApiKey})
			ctx := metadata.NewOutgoingContext(context.Background(), header)

			resp, err := c.DeleteStream(ctx, req)
			if err != nil {
				if tt.ignoreErr == "" || (tt.ignoreErr != "" && !strings.Contains(err.Error(), tt.ignoreErr)) {
					errStatus, _ := status.FromError(err)
					t.Errorf("%s got unexpected error %q", tt.name, errStatus.Message())
				}
				break
			}
			if resp.Id != tt.want {
				t.Errorf("%v, \nwant contain %q, \nhave %q", tt.name, tt.want, resp.Id)
			}

		}

	}
}
