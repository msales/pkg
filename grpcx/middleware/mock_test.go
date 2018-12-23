package middleware_test

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type serverStreamMock struct {
	ctx context.Context
}

func (*serverStreamMock) SetHeader(metadata.MD) error {
	return nil
}

func (*serverStreamMock) SendHeader(metadata.MD) error {
	return nil
}

func (*serverStreamMock) SetTrailer(metadata.MD) {
}

func (s *serverStreamMock) Context() context.Context {
	return s.ctx
}

func (*serverStreamMock) SendMsg(m interface{}) error {
	return nil
}

func (*serverStreamMock) RecvMsg(m interface{}) error {
	return nil
}
