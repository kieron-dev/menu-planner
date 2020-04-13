package routing_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/kieron-pivotal/menu-planner-app/routing"
	"github.com/kieron-pivotal/menu-planner-app/routing/routingfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {

	Context("routing", func() {
		var (
			mockServer  *httptest.Server
			handlers    *routingfakes.FakeHandlers
			frontendURI = "https://foo.com"
		)

		BeforeEach(func() {
			handlers = new(routingfakes.FakeHandlers)
			router := routing.New(frontendURI, handlers)
			mockServer = httptest.NewServer(router.SetupRoutes())
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
		})

		It("calls authGoogle handler on POST", func() {
			_, err := http.Post(mockServer.URL+"/authGoogle", "application/json", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(handlers.AuthGoogleCallCount()).To(Equal(1))
		})
	})

})
