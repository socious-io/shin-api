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

func verificationGroup() {

	It("it should create verification", func() {
		for i, data := range verificationsData {
			w := httptest.NewRecorder()
			data["schema_id"] = schemasData[0]["id"]

			for j, attr := range schemasData[0]["attributes"].([]gin.H) {
				data["attributes"].([]gin.H)[j]["attribute_id"] = attr["id"]
			}

			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/verifications", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			verificationsData[i]["id"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should update verification", func() {
		for _, data := range verificationsData {
			w := httptest.NewRecorder()
			data["schema_id"] = schemasData[0]["id"]
			fmt.Println()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/verifications/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(202))
		}
	})

	It("it should get verification", func() {
		for _, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/verifications/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	})

	It("it should create verification individual", func() {
		for i, data := range verificationsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(gin.H{
				"customer_id":     fmt.Sprintf("customer_%d", i),
				"verification_id": data["id"],
			})
			req, _ := http.NewRequest("POST", "/verifications/individuals", bytes.NewBuffer(reqBody))
			// req.Header.Set("apikey", intKey)
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			verificationsData[i]["individual"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should get verification individual by customer", func() {
		for i, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/verifications/%s/individuals/customer_%d", data["id"], i), nil)
			req.Header.Set("apikey", intKey)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
		}
	})

	/* It("it should get verification with connection", func() {
		for _, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/verifications/%s/connect", data["individual"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["individual"]).To(Equal(data["id"]))
			// Expect(body["connection_url"]).To(Not(nil))
			Expect(w.Code).To(Equal(200))
		}
	}) */

	It("it should get verirications", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/verifications", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(int(body["total"].(float64))).To(Equal(len(verificationsData)))
		Expect(w.Code).To(Equal(200))
	})

	It("it should delete verification", func() {
		for _, data := range verificationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/verifications/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(200))
		}
	})

}
