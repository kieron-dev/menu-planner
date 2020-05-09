package handlers_test

import (
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Cookies and sessions", func() {
	Context("session fixation", func() {
		When("an invalid session cookie is received", func() {
			It("deletes cookie and invalid input", func() {

			})
		})

		When("auth is successful", func() {
			It("generates a new session id", func() {

			})
		})
	})

	When("the user logs out", func() {
		It("clears the session cookie", func() {

		})

		It("deletes the session", func() {

		})
	})

})
