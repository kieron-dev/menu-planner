package db_test

import (
	"database/sql"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	_ "github.com/lib/pq"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var (
	pg *sql.DB
)

var _ = BeforeSuite(func() {
	connStr := mustGetEnv("DB_CONN_STR")
	var err error
	pg, err = sql.Open("postgres", connStr)
	Expect(err).NotTo(HaveOccurred())
})

func mustGetEnv(v string) string {
	s := os.Getenv(v)
	if s != "" {
		return s
	}
	Fail(fmt.Sprintf("expected env var %q", v))
	return ""
}
