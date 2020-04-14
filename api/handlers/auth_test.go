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
		httpHandlers  *handlers.Handlers
		hf            http.HandlerFunc
		tokenVerifier *handlersfakes.FakeTokenVerifier
		jwtDecoder    *handlersfakes.FakeJWTDecoder
		localAuther   *handlersfakes.FakeLocalAuther
		recorder      *httptest.ResponseRecorder
		req           *http.Request
		audience      string
		bodyBytes     []byte
	)

	BeforeEach(func() {
		log.SetOutput(GinkgoWriter)
		tokenVerifier = new(handlersfakes.FakeTokenVerifier)
		jwtDecoder = new(handlersfakes.FakeJWTDecoder)
		jwtDecoder.ClaimSetReturns(map[string]string{"email": "bar@foo.com", "name": "bar"}, nil)
		localAuther = new(handlersfakes.FakeLocalAuther)
		audience = "my.audience"
		httpHandlers = handlers.New(audience, tokenVerifier, jwtDecoder, localAuther)
		hf = http.HandlerFunc(httpHandlers.AuthGoogle)
		recorder = httptest.NewRecorder()
		bodyBytes = []byte("{}")
	})

	JustBeforeEach(func() {
		body := bytes.NewBuffer(bodyBytes)
		var err error
		req, err = http.NewRequest(http.MethodPost, "application/json", body)
		Expect(err).NotTo(HaveOccurred())
		hf.ServeHTTP(recorder, req)
	})

	Context("google auth", func() {

		When("the token is valid", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`{"tokenID":"my.google.token"}`)
				jwtDecoder.ClaimSetReturns(map[string]string{"name": "bob", "email": "bob@bits.com"}, nil)
				localAuther.LocalAuthReturns("a.valid.token", nil)
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

			It("succeeds", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(Equal(`{"token":"a.valid.token"}`))
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
				jwtDecoder.ClaimSetReturns(map[string]string{"foo": "bar"}, nil)
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("the token doesn't include name", func() {
			BeforeEach(func() {
				jwtDecoder.ClaimSetReturns(map[string]string{"email": "bar@foo.com"}, nil)
			})

			It("fails with bad request error", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("local auth fails", func() {
			BeforeEach(func() {
				localAuther.LocalAuthReturns("", errors.New("oops"))
			})

			It("returns an internal server error", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

	})
})
