package integration_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	BeforeEach(func() {
		frontendURI = "https://my.frontend.com"
		tokenVerifier = new(handlersfakes.FakeTokenVerifier)
	})

	JustBeforeEach(func() {
		h := handlers.NewAuthHandler(audience, tokenVerifier, jwtDecoder, userStore, sessionManager)
		r := routing.New(frontendURI, sessionManager, h)
		mockServer = httptest.NewServer(r.SetupRoutes())
	})

	Context("auth", func() {
		It("cannot access /whoami un-authed", func() {
			resp, err := http.Get(mockServer.URL + "/whoami")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		When("it receives a valid google JWT", func() {
			BeforeEach(func() {
				tokenVerifier.VerifyIDTokenReturns(nil)
			})

			It("returns a session cookie which can access privileged routes", func() {
				jstr := `{"email":"foo@bar.com", "name":"foo bar"}`
				b64str := base64.StdEncoding.EncodeToString([]byte(jstr))
				data := fmt.Sprintf(`{"tokenID": "xxx.%s.zzz"}`, b64str)

				resp, err := http.Post(mockServer.URL+"/authGoogle", "application/json", bytes.NewBufferString(data))
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				defer resp.Body.Close()

				cookies := resp.Cookies()
				Expect(cookies).To(HaveLen(1))

				req, err := http.NewRequest(http.MethodGet, mockServer.URL+"/whoami", nil)
				Expect(err).NotTo(HaveOccurred())
				req.AddCookie(cookies[0])

				resp, err = http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(body)).To(ContainSubstring("Hello, foo bar"))
			})
		})
	})
})
