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

func credentialGroup() {
	if focused {
		authGroup()
		orgGroup()
		recipientGroup()
		schemaGroup()
	}

	It("it should create credential", func() {
		for i, data := range credentialsData {
			w := httptest.NewRecorder()
			data["recipient_id"] = recipientsData[0]["id"]
			data["schema_id"] = schemasData[0]["id"]
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/credentials", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			credentialsData[i]["id"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should update credential", func() {
		for i, data := range credentialsData {
			w := httptest.NewRecorder()
			name := fmt.Sprintf("test_name_%d", i+1)
			data["name"] = name
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/credentials/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["name"].(string)).To(Equal(name))
			Expect(w.Code).To(Equal(202))
		}
	})

	It("it should get credential", func() {
		for _, data := range credentialsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	})

	/* It("it should get credential with connection", func() {
		for _, data := range credentialsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/%s/connect", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	}) */

	It("it should get credentials", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/credentials", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(len(body["results"].([]interface{}))).To(Equal(len(credentialsData)))
		Expect(w.Code).To(Equal(200))
	})

	It("it should get credentials by schema", func() {
		w := httptest.NewRecorder()
		schemaId := credentialsData[0]["schema_id"]
		req, _ := http.NewRequest("GET", fmt.Sprintf("/credentials?schema=%s", schemaId), nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(len(body["results"].([]interface{}))).To(Equal(len(credentialsData)))
		Expect(w.Code).To(Equal(200))
	})

	It("it should delete credential", func() {
		for _, data := range credentialsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/credentials/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
		}
	})

	It("it should create recipient+credential", func() {
		for i, cdata := range credentialsData {
			w := httptest.NewRecorder()
			cdata["schema_id"] = schemasData[0]["id"]
			data := map[string]any{
				"recipient":  recipientsData[0],
				"credential": cdata,
			}
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/credentials/with-recipient", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			credentialsData[i]["id"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})
}
