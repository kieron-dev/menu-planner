package db_test

import (
	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recipe", func() {
	var recipeStore *db.RecipeStore

	Describe("Listing recipes", func() {
		var (
			recipes []models.Recipe
			err     error
			userID  int
		)

		BeforeEach(func() {
			userID = 234
			recipeStore = db.NewRecipeStore(tx)

			_, err := tx.Exec(`insert into local_user(id, name, email)
                    VALUES (123, 'bob', 'bob@example.com'),
                    (234, 'jim', 'jim@example.com'),
                    (345, 'gertrude', 'gertrude@example.com')`)
			Expect(err).NotTo(HaveOccurred())

			_, err = tx.Exec(`insert into recipe (name, user_id)
               VALUES ('recipe 1', 123),
               ('recipe 2', 234),
               ('recipe 3', 234),
               ('recipe 4', 345)`)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			recipes, err = recipeStore.List(userID)
		})

		When("there are no recipes for the user", func() {
			BeforeEach(func() {
				userID = 1
			})

			It("returns an empty slice", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(recipes).To(BeEmpty())
			})
		})

		When("there are some recipes with matching user id", func() {
			It("returns an recipes 2 and 3", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(recipes).To(HaveLen(2))
				Expect(recipes[0].Name()).To(Equal("recipe 2"))
				Expect(recipes[1].Name()).To(Equal("recipe 3"))
			})
		})
	})
})
