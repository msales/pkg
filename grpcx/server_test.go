package grpcx_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/msales/pkg/v4/grpcx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	srv := grpcx.NewServer(new(mockListener), &grpc.Server{})

	assert.IsType(t, (*grpcx.Server)(nil), srv)
}

func TestServer_Serve(t *testing.T) {
	lisErr := errors.New("test: listener error")

	lis := new(mockListener)
	lis.On("Accept").Return(new(mockConn), lisErr)
	lis.On("Close").Return(nil)

	srv := grpcx.NewServer(lis, grpc.NewServer())

	err := srv.Serve()

	assert.Equal(t, lisErr, err)
	lis.AssertExpectations(t)
}

func TestServer_Close(t *testing.T) {
	// Cannot properly unit-test it because the server cannot be mocked.
}

type mockListener struct {
	mock.Mock
}

func (l *mockListener) Accept() (net.Conn, error) {
	args := l.Called()

	return args.Get(0).(net.Conn), args.Error(1)
}

func (l *mockListener) Close() error {
	return l.Called().Error(0)
}

func (l *mockListener) Addr() net.Addr {
	return l.Called().Get(0).(net.Addr)
}

type mockConn struct{}

func (mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (mockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (mockConn) Close() error {
	return nil
}

func (mockConn) LocalAddr() net.Addr {
	return nil
}

func (mockConn) RemoteAddr() net.Addr {
	return nil
}

func (mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
