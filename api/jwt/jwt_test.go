package jwt_test

import (
	"encoding/base64"
	"fmt"

	"github.com/kieron-pivotal/menu-planner-app/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jwt", func() {

	var (
		token    string
		j        *jwt.JWT
		claimSet map[string]interface{}
		err      error
	)

	BeforeEach(func() {
		j = jwt.NewJWT()
	})

	JustBeforeEach(func() {
		claimSet, err = j.ClaimSet(token)
	})

	When("token is valid", func() {
		plain := `{"foo": "bar", "valid": true, "count": 42}`
		b64 := base64.StdEncoding.EncodeToString([]byte(plain))

		BeforeEach(func() {
			token = fmt.Sprintf("xxx.%s.yyy", b64)
		})

		It("successfully decodes the middle part", func() {
			foo, ok := claimSet["foo"]
			Expect(ok).To(BeTrue(), "foo wasn't present in claimSet")
			Expect(foo).To(Equal("bar"))

			valid, ok := claimSet["valid"]
			Expect(ok).To(BeTrue(), "valid wasn't present in claimSet")
			Expect(valid).To(BeTrue())

			count, ok := claimSet["count"]
			Expect(ok).To(BeTrue(), "count wasn't present in claimSet")
			Expect(count).To(Equal(float64(42)))

		})
	})

	When("token does not have 3 dot separated parts", func() {
		BeforeEach(func() {
			token = "asdf.asdf.fdsa.fdsa"
		})

		It("fails", func() {
			Expect(err).To(MatchError(ContainSubstring("invalid-format")))
		})
	})

	When("token is not base64", func() {
		BeforeEach(func() {
			token = "xxx.\x01\x02.yyy"
		})

		It("fails", func() {
			Expect(err).To(MatchError(ContainSubstring("decoding-token-failed")))
		})
	})

	When("claimSet isn't valid json map", func() {
		BeforeEach(func() {
			plain := "oh hello there"
			b64 := base64.StdEncoding.EncodeToString([]byte(plain))
			token = fmt.Sprintf("xxx.%s.yyy", b64)
		})

		It("fails", func() {
			Expect(err).To(MatchError(ContainSubstring("unmarshalling-token-failed")))
		})
	})

})
