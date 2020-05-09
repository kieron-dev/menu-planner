package jwt_test

import (
	"fmt"

	"github.com/kieron-pivotal/menu-planner-app/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokengenerator", func() {
	var (
		tokengen *jwt.JWT
		token    string
		id       string
		name     string
		claimSet map[string]interface{}
		err      error
	)

	BeforeEach(func() {
		tokengen = jwt.NewJWT()
	})

	JustBeforeEach(func() {
		token, err = tokengen.GenerateToken(id, name)
	})

	When("all goes smoothly", func() {
		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			claimSet, err = tokengen.ClaimSet(token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("has id in the claimSet", func() {
			hasKeyVal(claimSet, "id", id)
		})

	})
})

func hasKeyVal(claimSet map[string]interface{}, key string, val interface{}) {
	v, ok := claimSet[key]
	Expect(ok).To(BeTrue(), fmt.Sprintf("%q not present in claimSet", key))
	Expect(v).To(Equal(val))
}
