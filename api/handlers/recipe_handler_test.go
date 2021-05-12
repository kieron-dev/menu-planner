package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/models"
	"github.com/kieron-pivotal/menu-planner-app/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RecipeHandler", func() {
	var (
		sessionManager *handlersfakes.FakeSessionManager
		recipeStore    *handlersfakes.FakeRecipeStore
		recorder       *httptest.ResponseRecorder
		req            *http.Request
		httpHandlers   *handlers.RecipeHandler
		hf             http.HandlerFunc
		recipe1        models.Recipe
		recipe2        models.Recipe
	)

	BeforeEach(func() {
		sessionManager = new(handlersfakes.FakeSessionManager)
		recipeStore = new(handlersfakes.FakeRecipeStore)
		httpHandlers = handlers.NewRecipeHandler(sessionManager, recipeStore)
		recorder = httptest.NewRecorder()
		recipe1 = models.Recipe{Name: "Bob", ID: 345}
		recipe2 = models.Recipe{Name: "Jim", ID: 456}
	})

	Describe("GetRecipes", func() {
		BeforeEach(func() {
			hf = http.HandlerFunc(httpHandlers.GetRecipes)
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
				sessionManager.GetReturns(&session.AuthInfo{IsLoggedIn: false}, nil)
			})

			It("returns a status not auth'ed", func() {
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		When("I'm logged in", func() {
			BeforeEach(func() {
				sessionManager.GetReturns(&session.AuthInfo{IsLoggedIn: true, Name: "forest", ID: 234}, nil)
				recipeStore.ListReturns([]models.Recipe{recipe1, recipe2}, nil)
			})

			It("returns a status OK", func() {
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
			})

			It("lists recipes using user ID", func() {
				Expect(recipeStore.ListCallCount()).To(Equal(1))
				userID := recipeStore.ListArgsForCall(0)
				Expect(userID).To(Equal(234))
			})

			It("formats the returned recipes as JSON", func() {
				Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
				Expect(recorder.Body.String()).To(ContainSubstring(`[{"name":"Bob","id":345},{"name":"Jim","id":456}]`))
			})
		})
	})

	Context("NewRecipe", func() {
		var body io.Reader

		BeforeEach(func() {
			body = strings.NewReader("")
			hf = http.HandlerFunc(httpHandlers.NewRecipe)
		})

		JustBeforeEach(func() {
			var err error
			req, err = http.NewRequest(http.MethodPost, "application/json", body)
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
				sessionManager.GetReturns(&session.AuthInfo{IsLoggedIn: false}, nil)
			})

			It("returns a status not auth'ed", func() {
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		When("a json body with a recipe name is passed", func() {
			BeforeEach(func() {
				sessionManager.GetReturns(&session.AuthInfo{IsLoggedIn: true, Name: "forest", ID: 234}, nil)
				body = strings.NewReader(`{"name":"foo bar"}`)
				recipeStore.InsertReturns(models.Recipe{Name: "foo bar", ID: 456}, nil)
			})

			It("inserts the meal into the database", func() {
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusCreated))
				Expect(recipeStore.InsertCallCount()).To(Equal(1))
				recipe := recipeStore.InsertArgsForCall(0)
				Expect(recipe).To(Equal(models.Recipe{Name: "foo bar", ID: 0, UserID: 234}))
				Expect(recorder.Body.String()).To(SatisfyAll(
					ContainSubstring(`"name":"foo bar"`),
					ContainSubstring(`"id":456`),
				))
			})
		})
	})
})
