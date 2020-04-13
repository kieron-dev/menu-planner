package integration_test

import (
	"net/http/httptest"
	"testing"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	"github.com/kieron-pivotal/menu-planner-app/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	frontendURI   string
	tokenVerifier *handlersfakes.FakeTokenVerifier
	mockServer    *httptest.Server
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	localAuther := auth.New()
	jwtDecoder := jwt.NewJWT()
	h := handlers.New(tokenVerifier, jwtDecoder, localAuther)
	r := routing.New(frontendURI, h)
	mockServer = httptest.NewServer(r.SetupRoutes())
})

var _ = BeforeEach(func() {
	frontendURI = "https://my.frontend.com"
	tokenVerifier = new(handlersfakes.FakeTokenVerifier)
})
