package db_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kieron-pivotal/menu-planner-app/db"
)

var _ = Describe("User", func() {
	var (
		store *db.UserStore
		user  db.User
		err   error
		email string
		name  string
		uuid  []uint8
	)

	BeforeEach(func() {
		store = db.NewUserStore(pg)
		email = "jill@example.com"
		name = "Jillian Ex"
	})

	AfterEach(func() {
		_, err := pg.Exec(`
DELETE FROM local_user
WHERE email = $1`, email)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("FindByEmail", func() {
		JustBeforeEach(func() {
			user, err = store.FindByEmail(email)
		})

		When("a user with the email exists in the DB", func() {
			BeforeEach(func() {
				err := pg.QueryRow(`
INSERT INTO local_user (email, name)
VALUES ($1, $2)
RETURNING lid`, email, name).Scan(&uuid)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Email()).To(Equal(email))
				Expect(user.Id()).To(Equal(uuid))
			})
		})

		When("a user with the email does not exist", func() {
			It("returns a not-found error", func() {
				Expect(db.IsNotFoundErr(err)).To(BeTrue())
			})
		})
	})

	Context("Create", func() {
		JustBeforeEach(func() {
			user, err = store.Create(email, name)
		})

		When("all goes well", func() {
			It("creates a new user", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Email()).To(Equal(email))
				Expect(user.Name()).To(Equal(name))
				Expect(user.Id()).ToNot(BeEmpty())
			})
		})

		When("the sql insert fails", func() {
			BeforeEach(func() {
				// trigger value too long error
				name = strings.Repeat("A", 256)
			})

			It("returns an error", func() {
				Expect(err).To(MatchError(ContainSubstring("create-user failed")))
			})
		})
	})
})
