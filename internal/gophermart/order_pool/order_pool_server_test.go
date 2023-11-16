package orderpool

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type AccrualServer struct {
	srv  *httptest.Server
	port int

	successRequest bool
	T              *testing.T
}

func NewAccrualServer(t *testing.T) *AccrualServer {
	accrual := &AccrualServer{T: t}
	accrual.srv = httptest.NewUnstartedServer(http.HandlerFunc(accrual.handlerSuccess))

	return accrual
}

func (srv *AccrualServer) Run() error {
	if err := srv.srv.Listener.Close(); err != nil {
		return err
	}

	if err := srv.getFreePort(); err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.port))
	if err != nil {
		return err
	}

	srv.srv.Listener = l
	srv.srv.Start()

	return nil
}

func (srv *AccrualServer) Close() {
	srv.srv.Close()
}

func (srv *AccrualServer) GetAddress() string {
	return fmt.Sprintf("http://localhost:%d", srv.port)
}

func (srv *AccrualServer) SetSuccessRequest(v bool) {
	srv.successRequest = v
}

func (srv *AccrualServer) handlerSuccess(w http.ResponseWriter, _ *http.Request) {
	if srv.successRequest {
		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(`{"order":"123","status":"PROCESSED","accrual":35.2}`))
		require.NoError(srv.T, err, "Не удалось записать ответ от тестового сервера accrual")
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (srv *AccrualServer) getFreePort() error {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	srv.port = l.Addr().(*net.TCPAddr).Port
	return nil
}
