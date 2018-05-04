package server_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpordomingo/gitbase-playground/server"
	"github.com/dpordomingo/gitbase-playground/server/handler"
	"github.com/dpordomingo/gitbase-playground/server/service"

	"github.com/stretchr/testify/suite"
)

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

type RouterTestSuite struct {
	ClientTestSuite
	router http.Handler
	server *httptest.Server
	db     service.SQLDB
}

const version = "test-version"

func (s *RouterTestSuite) SetupSuite() {
	logger := service.NewLogger("dev")
	staticHandler := handler.NewStatic("/tmp", "localhost:0")
	s.db = &mockDB{}
	s.router = server.Router(
		logger,
		staticHandler,
		version,
		s.db,
	)
}

func (s *RouterTestSuite) SetupTest() {
	s.server = httptest.NewServer(s.router)
}

func (s *RouterTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *RouterTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *RouterTestSuite) GetResponse(method string, path string, body io.Reader) *http.Response {
	url := s.server.URL + path
	response, err := GetResponse(method, url, body)
	if err != nil {
		s.Fail(err.Error())
	}

	return response
}

func (s *RouterTestSuite) TestVersion() {
	expectedVersion := fmt.Sprintf(`{"status":200,"data":{"version":"%s"}}`, version)
	response := s.GetResponse("GET", "/version", nil)
	s.AssertResponseBodyStatus(response, 200, expectedVersion, "version should be served")
}
