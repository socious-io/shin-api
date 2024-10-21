package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func recipientGroup() {

	It("it should create recipient", func() {
		for i, data := range recipientsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/recipients", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			recipientsData[i]["id"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should update recipient", func() {
		for i, data := range recipientsData {
			w := httptest.NewRecorder()
			lastName := fmt.Sprintf("test_last_%d", i+1)
			data["last_name"] = lastName
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/recipients/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["last_name"].(string)).To(Equal(lastName))
			Expect(w.Code).To(Equal(202))
		}
	})

	It("it should get recipient", func() {
		for _, data := range recipientsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/recipients/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	})

	It("it should get recipients", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipients", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(len(body["results"].([]interface{}))).To(Equal(len(recipientsData)))
		Expect(w.Code).To(Equal(200))
	})

	It("it should delete recipients", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/recipients/%s", recipientsData[2]["id"]), nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		fmt.Println(body)
		// Expect(body["id"]).To(Equal(data["id"]))
		Expect(w.Code).To(Equal(200))
	})

}
