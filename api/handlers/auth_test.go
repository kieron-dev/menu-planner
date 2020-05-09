package handlers_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"

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
		localAuther   *handlersfakes.FakeLocalAuther
		sessionSetter *handlersfakes.FakeSessionSetter
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
		localAuther = new(handlersfakes.FakeLocalAuther)
		user = new(handlersfakes.FakeUser)
		audience = "my.audience"
		sessionSetter = new(handlersfakes.FakeSessionSetter)
		httpHandlers = handlers.New(audience, tokenVerifier, jwtDecoder, localAuther, sessionSetter)
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
				localAuther.LocalAuthReturns(user, nil)
				user.IDReturns(12345)
				user.NameReturns("user-name")
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

			It("sends email and name to local auth'er", func() {
				Expect(localAuther.LocalAuthCallCount()).To(Equal(1))
				email, name := localAuther.LocalAuthArgsForCall(0)
				Expect(email).To(Equal("bob@bits.com"))
				Expect(name).To(Equal("bob"))
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

		When("the token doesn't include name", func() {
			BeforeEach(func() {
				jwtDecoder.ClaimSetReturns(map[string]interface{}{"email": "bar@foo.com"}, nil)
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("local auth fails", func() {
			BeforeEach(func() {
				localAuther.LocalAuthReturns(nil, errors.New("oops"))
			})

			It("returns an internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

	})
})
