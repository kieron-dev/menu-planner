package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/handlers/handlersfakes"
	"github.com/kieron-pivotal/menu-planner-app/models"
	"github.com/kieron-pivotal/menu-planner-app/models/modelsfakes"
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
		recipe1        *modelsfakes.FakeRecipe
		recipe2        *modelsfakes.FakeRecipe
	)

	BeforeEach(func() {
		sessionManager = new(handlersfakes.FakeSessionManager)
		recipeStore = new(handlersfakes.FakeRecipeStore)
		httpHandlers = handlers.NewRecipeHandler(sessionManager, recipeStore)
		hf = http.HandlerFunc(httpHandlers.GetRecipes)
		recorder = httptest.NewRecorder()
		recipe1 = new(modelsfakes.FakeRecipe)
		recipe2 = new(modelsfakes.FakeRecipe)
	})

	JustBeforeEach(func() {
		var err error
		req, err = http.NewRequest(http.MethodGet, "application/json", nil)
		Expect(err).NotTo(HaveOccurred())
		hf.ServeHTTP(recorder, req)
	})

	Describe("GetRecipes", func() {
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
				recipe1.NameReturns("Bob")
				recipe1.IDReturns(345)
				recipe2.NameReturns("Jim")
				recipe2.IDReturns(456)
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
				Expect(string(recorder.Body.Bytes())).To(ContainSubstring(`[{"name":"Bob","id":345},{"name":"Jim","id":456}]`))
			})
		})
	})
})
