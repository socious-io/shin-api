package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func integrationGroup() {

	It("it should create integration", func() {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(&gin.H{"name": "Test"})
		req, _ := http.NewRequest("POST", "/integrations", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		intKey = body["secret"].(string)
		Expect(w.Code).To(Equal(201))
	})
}
