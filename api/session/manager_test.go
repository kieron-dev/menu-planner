package session_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/securecookie"
	"github.com/kieron-pivotal/menu-planner-app/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var (
		sessionManager    *session.Manager
		sessionKeys       [][]byte
		middleware        http.Handler
		next              http.Handler
		req               *http.Request
		resp              *httptest.ResponseRecorder
		sessionCookieName = "_id"
		err               error
		lambda            func(w http.ResponseWriter, r *http.Request)
	)

	BeforeEach(func() {
		lambda = func(w http.ResponseWriter, r *http.Request) {}
		sessionKeys = [][]byte{securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)}
		sessionManager = session.NewManager(sessionKeys)
		req, err = http.NewRequest(http.MethodGet, "", nil)
		Expect(err).NotTo(HaveOccurred())
		resp = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lambda(w, r)
		})
		middleware = sessionManager.SessionMiddleware(next)
		middleware.ServeHTTP(resp, req)
	})

	When("invoking the 'next' middleware", func() {
		var called bool

		BeforeEach(func() {
			lambda = func(w http.ResponseWriter, r *http.Request) {
				called = true
			}
		})

		It("does actually call it", func() {
			Expect(called).To(BeTrue())
		})
	})

	When("a session cookie is invalid", func() {
		BeforeEach(func() {
			sessionCookie := &http.Cookie{
				Name:  sessionCookieName,
				Value: "definitely-not-valid",
			}
			req.AddCookie(sessionCookie)
		})

		It("rejects invalid sessions with a invalid request error", func() {
			Expect(resp.Result().StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("deletes that cookie", func() {
			cookies := resp.Result().Cookies()
			Expect(cookies).To(HaveLen(1))
			cookie := cookies[0]
			Expect(cookie.Name).To(Equal(sessionCookieName))
			Expect(cookie.MaxAge).To(BeNumerically("<", 0))
		})
	})

	When("no session cookie is passed in, but session is saved", func() {
		BeforeEach(func() {
			lambda = func(w http.ResponseWriter, r *http.Request) {
				sessionManager.Set(r, w, &session.AuthInfo{Name: "alice"})
			}
		})

		It("creates a new session cookie", func() {
			cookies := resp.Result().Cookies()
			Expect(cookies).To(HaveLen(1))
			cookie := cookies[0]
			Expect(cookie.Name).To(Equal(sessionCookieName))
		})
	})

	Context("session storage", func() {
		var (
			ourSession *session.AuthInfo
			setErr     error
		)

		BeforeEach(func() {
			ourSession = &session.AuthInfo{
				Name:       "bob",
				ID:         10,
				IsLoggedIn: true,
			}
			lambda = func(w http.ResponseWriter, r *http.Request) {
				setErr = sessionManager.Set(r, w, ourSession)
			}
		})

		It("can store things in the session", func() {
			Expect(setErr).NotTo(HaveOccurred())
		})

		When("performing a subsequent request", func() {
			It("can get the session inside the handlers", func() {
				cookies := resp.Result().Cookies()
				Expect(cookies).To(HaveLen(1))

				req, err = http.NewRequest(http.MethodGet, "", nil)
				Expect(err).NotTo(HaveOccurred())
				req.AddCookie(cookies[0])
				resp = httptest.NewRecorder()

				var sess *session.AuthInfo
				next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var err error
					sess, err = sessionManager.Get(r.Context())
					Expect(err).NotTo(HaveOccurred())
				})
				middleware = sessionManager.SessionMiddleware(next)
				middleware.ServeHTTP(resp, req)

				Expect(sess.Name).To(Equal("bob"))
				Expect(sess.IsLoggedIn).To(BeTrue())
				Expect(sess.ID).To(Equal(10))
			})
		})
	})
})
