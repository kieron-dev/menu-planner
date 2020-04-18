package auth_test

import (
	"errors"

	"github.com/kieron-pivotal/menu-planner-app/auth"
	"github.com/kieron-pivotal/menu-planner-app/auth/authfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalAuth", func() {
	var (
		localAuth        *auth.LocalAuth
		fakeUserStore    *authfakes.FakeUserStore
		fakeJWTGenerator *authfakes.FakeJWTGenerator
		fakeUser         *authfakes.FakeUser
		email            string
		name             string
		token            string
		err              error
	)

	BeforeEach(func() {
		fakeUserStore = new(authfakes.FakeUserStore)
		fakeJWTGenerator = new(authfakes.FakeJWTGenerator)
		fakeJWTGenerator.GenerateTokenReturns("abc.123.456", nil)
		localAuth = auth.NewLocalAuth(fakeUserStore, fakeJWTGenerator)
		fakeUser = new(authfakes.FakeUser)
		fakeUser.IDReturns("abc123")
		fakeUser.NameReturns(name)
		fakeUserStore.FindByEmailReturns(fakeUser, nil)

		email = "bob@example.com"
		name = "Robert Ample"
	})

	JustBeforeEach(func() {
		token, err = localAuth.LocalAuth(email, name)
	})

	When("all succeeds", func() {
		It("doesn't error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("checks the db for a row with the email", func() {
			Expect(fakeUserStore.FindByEmailCallCount()).To(Equal(1))
			actualEmail := fakeUserStore.FindByEmailArgsForCall(0)
			Expect(actualEmail).To(Equal(email))
		})

		When("the user with that email exists already", func() {
			It("doesn't create a new user", func() {
				Expect(fakeUserStore.CreateCallCount()).To(Equal(0))
			})

			It("uses the existing ID in the token generation", func() {
				Expect(fakeJWTGenerator.GenerateTokenCallCount()).To(Equal(1))
				actualID, _ := fakeJWTGenerator.GenerateTokenArgsForCall(0)
				Expect(actualID).To(Equal("abc123"))
			})
		})

		When("there is no existing user with that email", func() {
			BeforeEach(func() {
				fakeUserStore.FindByEmailReturns(nil, nil)
				fakeUser.IDReturns("new234")
				fakeUserStore.CreateReturns(fakeUser, nil)
			})

			It("creates a new user", func() {
				Expect(fakeUserStore.CreateCallCount()).To(Equal(1))
				actualEmail, actualName := fakeUserStore.CreateArgsForCall(0)
				Expect(actualEmail).To(Equal(email))
				Expect(actualName).To(Equal(name))
			})

			It("uses the new ID in the token generation", func() {
				Expect(fakeJWTGenerator.GenerateTokenCallCount()).To(Equal(1))
				actualID, _ := fakeJWTGenerator.GenerateTokenArgsForCall(0)
				Expect(actualID).To(Equal("new234"))
			})
		})

		It("uses the correct name in the JWT generation", func() {
			Expect(fakeJWTGenerator.GenerateTokenCallCount()).To(Equal(1))
			_, actualName := fakeJWTGenerator.GenerateTokenArgsForCall(0)
			Expect(actualName).To(Equal(name))
		})

		It("returns the token", func() {
			Expect(token).To(Equal("abc.123.456"))
		})
	})

	Context("sad paths", func() {
		When("find by email fails", func() {
			BeforeEach(func() {
				fakeUserStore.FindByEmailReturns(nil, errors.New("eek"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError(ContainSubstring("eek")))
			})
		})

		When("create user fails", func() {
			BeforeEach(func() {
				fakeUserStore.FindByEmailReturns(nil, nil)
				fakeUserStore.CreateReturns(nil, errors.New("mmm"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError(ContainSubstring("mmm")))
			})
		})

		When("the token generation fails", func() {
			BeforeEach(func() {
				fakeJWTGenerator.GenerateTokenReturns("", errors.New("sorry"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError(ContainSubstring("sorry")))
			})
		})
	})
})
