package db_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/models"
)

var _ = Describe("User", func() {
	var (
		store *db.UserStore
		user  models.User
		err   error
		email string
		name  string
		id    int
	)

	BeforeEach(func() {
		store = db.NewUserStore(tx)
		email = "jill@example.com"
		name = "Jillian Ex"
	})

	Context("FindByEmail", func() {
		JustBeforeEach(func() {
			user, err = store.FindByEmail(email)
		})

		When("a user with the email exists in the DB", func() {
			BeforeEach(func() {
				err := tx.QueryRow(`
INSERT INTO local_user (email, name)
VALUES ($1, $2)
RETURNING id`, email, name).Scan(&id)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Email()).To(Equal(email))
				Expect(user.ID()).To(Equal(id))
			})
		})

		When("a user with the email does not exist", func() {
			It("returns a not-found error", func() {
				Expect(store.IsNotFoundErr(err)).To(BeTrue())
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
				Expect(user.ID()).ToNot(BeZero())
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
