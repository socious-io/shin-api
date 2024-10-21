package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func schemaGroup() {

	It("it should create schema", func() {
		for i, data := range schemasData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/schemas", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			schemasData[i]["id"] = body["id"]
			attrs := body["attributes"].([]interface{})
			for j, attr := range attrs {
				schemasData[i]["attributes"].([]gin.H)[j]["id"] = attr.(map[string]interface{})["id"]
			}
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should get schema", func() {
		for _, data := range schemasData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/schemas/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	})

	It("it should get schemas", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/schemas", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		// body := decodeBody(w.Body)
		// fmt.Println(body["results"])
		// Expect(body["id"]).To(Equal(data["id"]))
		Expect(w.Code).To(Equal(200))
	})

	It("it should delete schemas", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/schemas/%s", schemasData[1]["id"]), nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		fmt.Println(body)
		// Expect(body["id"]).To(Equal(data["id"]))
		Expect(w.Code).To(Equal(200))
	})

}
