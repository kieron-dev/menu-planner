package integration_test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/kieron-pivotal/menu-planner-app/auth"
	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	"github.com/kieron-pivotal/menu-planner-app/routing"
	"github.com/kieron-pivotal/menu-planner-app/session"
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
	connStr := mustGetEnv("DB_CONN_STR")
	pg, err := sql.Open("postgres", connStr)
	Expect(err).NotTo(HaveOccurred())

	userStore := db.NewUserStore(pg)
	localAuther := auth.NewLocalAuth(userStore)
	jwtDecoder := jwt.NewJWT()
	audience := "my-web-app-id"
	h := handlers.New(audience, tokenVerifier, jwtDecoder, localAuther)
	sessionManager := session.NewManager([][]byte{securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32)})
	r := routing.New(frontendURI, sessionManager, h)
	mockServer = httptest.NewServer(r.SetupRoutes())
})

var _ = BeforeEach(func() {
	frontendURI = "https://my.frontend.com"
	tokenVerifier = new(handlersfakes.FakeTokenVerifier)
})

func mustGetEnv(v string) string {
	s := os.Getenv(v)
	if s != "" {
		return s
	}
	Fail(fmt.Sprintf("expected env var %q", v))
	return ""
}
