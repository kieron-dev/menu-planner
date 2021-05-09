package integration_test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	"github.com/kieron-pivotal/menu-planner-app/session"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	frontendURI    string
	tokenVerifier  *handlersfakes.FakeTokenVerifier
	mockServer     *httptest.Server
	audience       string
	userStore      *db.UserStore
	recipeStore    *db.RecipeStore
	jwtDecoder     *jwt.JWT
	sessionManager *session.Manager
	pg             *sql.DB
	tx             *sql.Tx
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	audience = "my-web-app-id"
	connStr := mustGetEnv("DB_CONN_STR")

	var err error
	pg, err = sql.Open("postgres", connStr)
	Expect(err).NotTo(HaveOccurred())

	jwtDecoder = jwt.NewJWT()

	keys := [][]byte{securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32)}
	sessionManager = session.NewManager(keys)
})

var _ = BeforeEach(func() {
	var err error
	tx, err = pg.Begin()
	Expect(err).NotTo(HaveOccurred())

	userStore = db.NewUserStore(tx)
	recipeStore = db.NewRecipeStore(tx)
})

var _ = AfterEach(func() {
	Expect(tx.Rollback()).To(Succeed())
})

func mustGetEnv(v string) string {
	s := os.Getenv(v)
	if s != "" {
		return s
	}
	Fail(fmt.Sprintf("expected env var %q", v))
	return ""
}
