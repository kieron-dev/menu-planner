package routing_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/kieron-pivotal/menu-planner-app/routing"
	"github.com/kieron-pivotal/menu-planner-app/routing/routingfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	Context("routing", func() {
		var (
			mockServer     *httptest.Server
			authHandler    *routingfakes.FakeAuthHandler
			recipeHandler  *routingfakes.FakeRecipeHandler
			frontendURI    = "https://foo.com"
			sessionManager *routingfakes.FakeSessionManager
		)

		BeforeEach(func() {
			authHandler = new(routingfakes.FakeAuthHandler)
			recipeHandler = new(routingfakes.FakeRecipeHandler)
			sessionManager = new(routingfakes.FakeSessionManager)
			// noop middleware
			sessionManager.SessionMiddlewareStub = func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			}
			router := routing.New(frontendURI, sessionManager, authHandler, recipeHandler)
			mockServer = httptest.NewServer(router.SetupRoutes())
		})

		Context("session management", func() {
			var (
				req *http.Request
				err error
			)

			JustBeforeEach(func() {
				req, err = http.NewRequest(http.MethodPost, mockServer.URL+"/authGoogle", nil)
				Expect(err).NotTo(HaveOccurred())
				_, err = http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
			})

			It("is wired into the middleware", func() {
				Expect(sessionManager.SessionMiddlewareCallCount()).To(Equal(1))
			})
		})

		Context("authGoogle", func() {
			It("gets CORS right", func() {
				req, err := http.NewRequest(http.MethodOptions, mockServer.URL+"/authGoogle", nil)
				Expect(err).NotTo(HaveOccurred())

				resp, err := http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				allowedOrigin := resp.Header.Get("Access-Control-Allow-Origin")
				Expect(allowedOrigin).To(Equal(frontendURI))

				allowedMethods := resp.Header.Get("Access-Control-Allow-Methods")
				Expect(allowedMethods).To(Equal("POST,OPTIONS"))

				allowedHeaders := resp.Header.Get("Access-Control-Allow-Headers")
				Expect(allowedHeaders).To(Equal("Content-Type"))
			})

			It("calls authGoogle handler on POST", func() {
				_, err := http.Post(mockServer.URL+"/authGoogle", "application/json", nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(authHandler.AuthGoogleCallCount()).To(Equal(1))
			})
		})

		Context("recipes", func() {
			It("calls getRecipes handler on GET /recipes", func() {
				_, err := http.Get(mockServer.URL + "/recipes")
				Expect(err).NotTo(HaveOccurred())
				Expect(recipeHandler.GetRecipesCallCount()).To(Equal(1))
			})

			It("calls newRecipe handler on POST /recipes", func() {
				body := strings.NewReader("")
				_, err := http.Post(mockServer.URL+"/recipes", "application/json", body)
				Expect(err).NotTo(HaveOccurred())
				Expect(recipeHandler.NewRecipeCallCount()).To(Equal(1))
			})
		})
	})
})
