package auth_test

import (
	"github.com/kieron-pivotal/menu-planner-app/auth"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalAuth", func() {
	var (
		localAuth *auth.LocalAuth
		email     string
		name      string
		token     string
		err       error
	)

	BeforeEach(func() {
		localAuth = auth.NewLocalAuth()
	})

	JustBeforeEach(func() {
		token, err = localAuth.LocalAuth(email, name)
	})

	When("all succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a token", func() {
			Expect(token).NotTo(BeEmpty())
		})

		It("contains email and name in the claim set", func() {
			claimSet, err := jwt.NewJWT().ClaimSet(token)
			Expect(err).NotTo(HaveOccurred())

			actualEmail, ok := claimSet["email"]
			Expect(ok).To(BeTrue(), "email not present")
			Expect(actualEmail).To(Equal(email))

			actualName, ok := claimSet["name"]
			Expect(ok).To(BeTrue(), "name not present")
			Expect(actualName).To(Equal(name))
		})
	})
})
