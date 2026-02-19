package base

import (
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"
)

type BaseIntegrationTestSuite struct {
	suite.Suite
	HttpClient *http.Client
}

func (s *BaseIntegrationTestSuite) SetupSuite() {
	s.HttpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}
