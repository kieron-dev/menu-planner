package integration_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {

	Context("auth", func() {
		When("it receives a valid google JWT", func() {

			It("returns an app JWT", func() {
				jstr := `{"email":"foo@bar.com", "name":"foo bar"}`
				b64str := base64.StdEncoding.EncodeToString([]byte(jstr))
				data := fmt.Sprintf(`{"tokenID": "xxx.%s.zzz"}`, b64str)

				resp, err := http.Post(mockServer.URL+"/authGoogle", "application/json", bytes.NewBufferString(data))
				Expect(err).NotTo(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				var respObj struct {
					Token string `json:"token"`
				}
				err = json.Unmarshal(body, &respObj)
				Expect(err).NotTo(HaveOccurred())

				Expect(respObj.Token).NotTo(BeEmpty())
			})
		})
	})

})
