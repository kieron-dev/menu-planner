package handlers_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {

	var (
		httpHandlers  *handlers.AuthHandler
		hf            http.HandlerFunc
		tokenVerifier *handlersfakes.FakeTokenVerifier
		jwtDecoder    *handlersfakes.FakeJWTDecoder
		sessionSetter *handlersfakes.FakeSessionSetter
		userStore     *handlersfakes.FakeUserStore
		user          *handlersfakes.FakeUser
		recorder      *httptest.ResponseRecorder
		req           *http.Request
		audience      string
		bodyBytes     []byte
		sessionCookie *http.Cookie
	)

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		tokenVerifier = new(handlersfakes.FakeTokenVerifier)
		jwtDecoder = new(handlersfakes.FakeJWTDecoder)
		jwtDecoder.ClaimSetReturns(map[string]interface{}{"email": "bar@foo.com", "name": "bar"}, nil)

		user = new(handlersfakes.FakeUser)
		user.IDReturns(12345)
		user.NameReturns("user-name")

		userStore = new(handlersfakes.FakeUserStore)
		userStore.FindByEmailReturns(user, nil)

		audience = "my.audience"
		sessionSetter = new(handlersfakes.FakeSessionSetter)
		httpHandlers = handlers.New(audience, tokenVerifier, jwtDecoder, userStore, sessionSetter)
		hf = http.HandlerFunc(httpHandlers.AuthGoogle)
		recorder = httptest.NewRecorder()
		bodyBytes = []byte("{}")
	})

	JustBeforeEach(func() {
		body := bytes.NewBuffer(bodyBytes)
		var err error
		req, err = http.NewRequest(http.MethodPost, "application/json", body)
		if sessionCookie != nil {
			req.AddCookie(sessionCookie)
		}
		Expect(err).NotTo(HaveOccurred())
		hf.ServeHTTP(recorder, req)
	})

	Context("google auth", func() {
		When("the token is valid", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`{"tokenID":"my.google.token"}`)
				jwtDecoder.ClaimSetReturns(map[string]interface{}{"name": "bob", "email": "bob@bits.com"}, nil)
			})

			It("calls the validator with correct args", func() {
				Expect(tokenVerifier.VerifyIDTokenCallCount()).To(Equal(1))
				token, aud := tokenVerifier.VerifyIDTokenArgsForCall(0)
				Expect(token).To(Equal("my.google.token"))
				Expect(aud).To(ConsistOf("my.audience"))
			})

			It("sends the token to the decoder", func() {
				Expect(jwtDecoder.ClaimSetCallCount()).To(Equal(1))
				Expect(jwtDecoder.ClaimSetArgsForCall(0)).To(Equal("my.google.token"))
			})

			It("tries to find user by email", func() {
				Expect(userStore.FindByEmailCallCount()).To(Equal(1))
				Expect(userStore.FindByEmailArgsForCall(0)).To(Equal("bob@bits.com"))
			})

			When("the user doesn't exist", func() {
				BeforeEach(func() {
					userStore.FindByEmailReturns(nil, db.NotFoundErr())
					userStore.CreateReturns(user, nil)
				})

				It("creates the user", func() {
					Expect(userStore.CreateCallCount()).To(Equal(1))
					actualEmail, actualName := userStore.CreateArgsForCall(0)
					Expect(actualEmail).To(Equal("bob@bits.com"))
					Expect(actualName).To(Equal("bob"))
				})
			})

			It("sets a new logged-in session", func() {
				Expect(sessionSetter.SetCallCount()).To(Equal(1))
				_, _, sess := sessionSetter.SetArgsForCall(0)
				Expect(sess.ID).To(Equal(12345))
				Expect(sess.Name).To(Equal("user-name"))
				Expect(sess.IsLoggedIn).To(BeTrue())
			})

			It("succeeds", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})

		When("the body is not valid json", func() {
			BeforeEach(func() {
				bodyBytes = []byte("{")
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("token verification fails", func() {
			BeforeEach(func() {
				tokenVerifier.VerifyIDTokenReturns(errors.New("expired"))
			})

			It("fails with a bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("decoding the claim set fails", func() {
			BeforeEach(func() {
				jwtDecoder.ClaimSetReturns(nil, errors.New("whoops"))
			})

			It("fails with a bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("the token doesn't include email address", func() {
			BeforeEach(func() {
				jwtDecoder.ClaimSetReturns(map[string]interface{}{"foo": "bar"}, nil)
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("the user must be created but the token doesn't include name", func() {
			BeforeEach(func() {
				userStore.FindByEmailReturns(nil, db.NotFoundErr())
				jwtDecoder.ClaimSetReturns(map[string]interface{}{"email": "bar@foo.com"}, nil)
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("findByEmail fails", func() {
			BeforeEach(func() {
				userStore.FindByEmailReturns(nil, errors.New("oops"))
			})

			It("fails with internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("create user fails", func() {
			BeforeEach(func() {
				userStore.FindByEmailReturns(nil, db.NotFoundErr())
				userStore.CreateReturns(nil, errors.New("oops"))
			})

			It("fails with internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("setting the session fails", func() {
			BeforeEach(func() {
				sessionSetter.SetReturns(errors.New("oops"))
			})

			It("fails with internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

		})
	})
})
