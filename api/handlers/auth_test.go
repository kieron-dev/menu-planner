package handlers_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/models/modelsfakes"
	"github.com/kieron-pivotal/menu-planner-app/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var (
		httpHandlers   *handlers.AuthHandler
		hf             http.HandlerFunc
		tokenVerifier  *handlersfakes.FakeTokenVerifier
		jwtDecoder     *handlersfakes.FakeJWTDecoder
		sessionManager *handlersfakes.FakeSessionManager
		userStore      *handlersfakes.FakeUserStore
		user           *modelsfakes.FakeUser
		recorder       *httptest.ResponseRecorder
		req            *http.Request
		audience       string
		bodyBytes      []byte
		sessionCookie  *http.Cookie
	)

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		tokenVerifier = new(handlersfakes.FakeTokenVerifier)
		jwtDecoder = new(handlersfakes.FakeJWTDecoder)
		jwtDecoder.ClaimSetReturns(map[string]interface{}{"email": "bar@foo.com", "name": "bar"}, nil)

		user = new(modelsfakes.FakeUser)
		user.IDReturns(12345)
		user.NameReturns("user-name")

		userStore = new(handlersfakes.FakeUserStore)
		userStore.FindByEmailReturns(user, nil)

		audience = "my.audience"
		sessionManager = new(handlersfakes.FakeSessionManager)
		httpHandlers = handlers.NewAuthHandler(audience, tokenVerifier, jwtDecoder, userStore, sessionManager)
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
					userStore.FindByEmailReturns(nil, errors.New("oops"))
					userStore.IsNotFoundErrReturns(true)
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
				Expect(sessionManager.SetCallCount()).To(Equal(1))
				_, _, sess := sessionManager.SetArgsForCall(0)
				Expect(sess.ID).To(Equal(12345))
				Expect(sess.Name).To(Equal("user-name"))
				Expect(sess.IsLoggedIn).To(BeTrue())
			})

			It("returns an ok success status", func() {
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
				userStore.FindByEmailReturns(nil, errors.New("oops"))
				userStore.IsNotFoundErrReturns(true)
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
				sessionManager.SetReturns(errors.New("oops"))
			})

			It("fails with internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

var _ = Describe("Who Am I?", func() {
	var (
		httpHandlers   *handlers.AuthHandler
		hf             http.HandlerFunc
		recorder       *httptest.ResponseRecorder
		req            *http.Request
		sessionManager *handlersfakes.FakeSessionManager
	)

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		sessionManager = new(handlersfakes.FakeSessionManager)
		httpHandlers = handlers.NewAuthHandler("", nil, nil, nil, sessionManager)
		hf = http.HandlerFunc(httpHandlers.WhoAmI)
		recorder = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		var err error
		req, err = http.NewRequest(http.MethodGet, "application/json", nil)
		Expect(err).NotTo(HaveOccurred())
		hf.ServeHTTP(recorder, req)
	})

	When("there is no session", func() {
		BeforeEach(func() {
			sessionManager.GetReturns(nil, errors.New("no session"))
		})

		It("returns a status not auth'ed", func() {
			Expect(recorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("I'm logged out", func() {
		BeforeEach(func() {
			sessionManager.GetReturns(&session.Session{IsLoggedIn: false}, nil)
		})

		It("returns a status not auth'ed", func() {
			Expect(recorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("I'm logged in", func() {
		BeforeEach(func() {
			sessionManager.GetReturns(&session.Session{IsLoggedIn: true, Name: "forest"}, nil)
		})

		It("returns OK status and prints my name", func() {
			Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
			body, err := ioutil.ReadAll(recorder.Result().Body)
			Expect(err).NotTo(HaveOccurred())
			defer recorder.Result().Body.Close()
			Expect(string(body)).To(ContainSubstring("Hello, forest"))
		})
	})
})
