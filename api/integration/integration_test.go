package integration_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/models"
	"github.com/kieron-pivotal/menu-planner-app/routing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	BeforeEach(func() {
		frontendURI = "https://my.frontend.com"
		tokenVerifier = new(handlersfakes.FakeTokenVerifier)

		authHandler := handlers.NewAuthHandler(audience, tokenVerifier, jwtDecoder, userStore, sessionManager)
		recipeHandler := handlers.NewRecipeHandler(sessionManager, recipeStore)
		r := routing.New(frontendURI, sessionManager, authHandler, recipeHandler)
		mockServer = httptest.NewServer(r.SetupRoutes())
	})

	login := func() (*http.Response, error) {
		jstr := `{"email":"foo@bar.com", "name":"foo bar"}`
		b64str := base64.StdEncoding.EncodeToString([]byte(jstr))
		loginData := fmt.Sprintf(`{"idToken": "xxx.%s.zzz"}`, b64str)

		return http.Post(mockServer.URL+"/authGoogle", "application/json", bytes.NewBufferString(loginData))
	}

	Context("auth", func() {
		It("cannot access /whoami un-authed", func() {
			resp, err := http.Get(mockServer.URL + "/whoami")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		When("it receives a valid google JWT", func() {
			It("returns a session cookie which can access privileged routes", func() {
				resp, err := login()
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

				Expect(string(body)).To(ContainSubstring(`{"name": "foo bar"}`))
			})

			When("an auth'ed session receives a logout", func() {
				var cookies []*http.Cookie

				It("logs the session out", func() {
					resp, err := login()
					Expect(err).NotTo(HaveOccurred())

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					defer resp.Body.Close()

					cookies = resp.Cookies()
					Expect(cookies).To(HaveLen(1))

					req, err := http.NewRequest(http.MethodGet, mockServer.URL+"/whoami", nil)
					Expect(err).NotTo(HaveOccurred())
					req.AddCookie(cookies[0])

					resp, err = http.DefaultClient.Do(req)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					cookies = resp.Cookies()
					Expect(cookies).To(HaveLen(1))

					req, err = http.NewRequest(http.MethodPost, mockServer.URL+"/logout", nil)
					Expect(err).NotTo(HaveOccurred())
					req.AddCookie(cookies[0])

					resp, err = http.DefaultClient.Do(req)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					cookies = resp.Cookies()
					Expect(cookies).To(HaveLen(2))

					req, err = http.NewRequest(http.MethodGet, mockServer.URL+"/whoami", nil)
					Expect(err).NotTo(HaveOccurred())
					req.AddCookie(cookies[1])

					resp, err = http.DefaultClient.Do(req)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})
		})
	})

	Context("recipes", func() {
		var (
			req    *http.Request
			resp   *http.Response
			cookie *http.Cookie
		)

		AfterEach(func() {
			resp.Body.Close()
		})

		Context("GET /recipes", func() {
			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequest(http.MethodGet, mockServer.URL+"/recipes", nil)
				Expect(err).NotTo(HaveOccurred())

				if cookie != nil {
					req.AddCookie(cookie)
				}

				resp, err = http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
			})

			When("not auth'ed", func() {
				It("returns an unauthorized status", func() {
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			When("auth'ed", func() {
				BeforeEach(func() {
					r, err := login()
					Expect(err).NotTo(HaveOccurred())
					Expect(r.StatusCode).To(Equal(http.StatusOK))
					cookies := r.Cookies()
					Expect(cookies).To(HaveLen(1))
					cookie = cookies[0]
				})

				It("returns an empty JSON list of recipes", func() {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

					b, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					Expect(string(b)).To(HavePrefix("[]"))
				})
			})
		})

		Context("POST /recipes", func() {
			var body io.Reader

			BeforeEach(func() {
				body = strings.NewReader(`{"name":"Roast Beef"}`)
			})

			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequest(http.MethodPost, mockServer.URL+"/recipes", body)
				Expect(err).NotTo(HaveOccurred())

				if cookie != nil {
					req.AddCookie(cookie)
				}

				resp, err = http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
			})

			When("not auth'ed", func() {
				It("returns an unauthorized status", func() {
					Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			When("auth'ed", func() {
				BeforeEach(func() {
					r, err := login()
					Expect(err).NotTo(HaveOccurred())
					Expect(r.StatusCode).To(Equal(http.StatusOK))
					cookies := r.Cookies()
					Expect(cookies).To(HaveLen(1))
					cookie = cookies[0]
				})

				It("returns the inserted recipe with ID", func() {
					Expect(resp.StatusCode).To(Equal(http.StatusCreated))
					Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

					b, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					var recipe models.Recipe
					err = json.Unmarshal(b, &recipe)
					Expect(err).NotTo(HaveOccurred())

					Expect(recipe.Name).To(Equal("Roast Beef"))
					Expect(recipe.ID).To(BeNumerically(">", 0))
				})
			})
		})
	})
})
