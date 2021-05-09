package db_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/lib/pq"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var (
	pg *sql.DB
	tx *sql.Tx
)

var _ = BeforeSuite(func() {
	connStr := mustGetEnv("DB_CONN_STR")
	var err error
	pg, err = sql.Open("postgres", connStr)
	Expect(err).NotTo(HaveOccurred())
})

var _ = BeforeEach(func() {
	var err error
	tx, err = pg.Begin()
	Expect(err).NotTo(HaveOccurred())
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
