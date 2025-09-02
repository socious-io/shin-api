package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"shin/src/config"
	"shin/src/shortner"
	"shin/src/utils"
	"strings"
	"time"
)

type H map[string]interface{}

type Connect struct {
	ID      string
	URL     string
	ShortID string
}

func CreateDID() (string, error) {

	// Document template with public keys and empty services
	documentTemplate := H{
		"publicKeys": []H{
			{
				"id":      "auth-1",
				"purpose": "authentication",
			},
			{
				"id":      "issue-1",
				"purpose": "assertionMethod",
			},
		},
		"services": []interface{}{},
	}

	// First API call to create DID
	didRes, err := makeRequest("/cloud-agent/did-registrar/dids", "POST", H{
		"documentTemplate": documentTemplate,
	})

	if err != nil {
		return "", err
	}

	var didResponse map[string]interface{}
	err = json.Unmarshal(didRes, &didResponse)
	if err != nil {
		return "", err
	}

	longFormDid := didResponse["longFormDid"].(string)

	// Second API call to publish DID
	pubRes, err := makeRequest(fmt.Sprintf("/cloud-agent/did-registrar/dids/%s/publications", longFormDid), "POST", H{})
	if err != nil {
		return "", err
	}

	var publishResponse H
	err = json.Unmarshal(pubRes, &publishResponse)
	if err != nil {
		return "", err
	}

	didRef := publishResponse["scheduledOperation"].(map[string]interface{})["didRef"].(string)
	return didRef, nil
}

func CreateConnection(callback string) (*Connect, error) {
	res, err := makeRequest("/cloud-agent/connections", "POST", H{"label": "Shin Connect"})
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}
	url := strings.ReplaceAll(
		body["invitation"].(map[string]interface{})["invitationUrl"].(string),
		"https://my.domain.com/path",
		config.Config.Wallet.Connect,
	)
	url += fmt.Sprintf("&callback=%s", callback)
	c := &Connect{
		ID:  body["connectionId"].(string),
		URL: url,
	}
	short, err := shortner.New(c.URL)
	if err != nil {
		return nil, err
	}
	c.ShortID = short.ShortID
	return c, nil
}

func ProofRequest(connectionID string, challenge string) (string, error) {
	time.Sleep(time.Second)
	res, err := makeRequest("/cloud-agent/present-proof/presentations", "POST", H{
		"connectionId": connectionID,
		"proofs":       []H{},
		"options": H{
			"challenge": challenge,
			"domain":    "shinid.com",
		},
	})
	if err != nil {
		return "", err
	}
	var body H

	if err := json.Unmarshal(res, &body); err != nil {
		return "", err
	}
	return body["presentationId"].(string), nil
}

func ProofVerify(presentID string) (H, error) {
	path := fmt.Sprintf("/cloud-agent/present-proof/presentations/%s", presentID)
	res, err := getRequest(path)
	if err != nil {
		return nil, err
	}
	var body H
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}

	if body["status"].(string) != "PresentationVerified" {
		return nil, fmt.Errorf("presentation not verified")
	}
	_, payload, err := utils.DecodeJWT(body["data"].([]interface{})[0].(string))
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode data")
	}
	var data H
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	_, payload, err = utils.DecodeJWT(data["vp"].(map[string]interface{})["verifiableCredential"].([]interface{})[0].(string))
	if err != nil {
		return nil, fmt.Errorf("presentation could not decode vc")
	}
	var vc H
	if err := json.Unmarshal(payload, &vc); err != nil {
		return nil, err
	}
	return vc["vc"].(map[string]interface{})["credentialSubject"].(map[string]interface{}), nil
}

func SendCredential(connectionID, did string, claims interface{}) (H, error) {
	payload := H{
		"claims":            claims,
		"connectionId":      connectionID,
		"issuingDID":        did,
		"schemaId":          nil,
		"automaticIssuance": true,
	}
	res, err := makeRequest("/cloud-agent/issue-credentials/credential-offers", "POST", payload)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	if err := json.Unmarshal(res, &body); err != nil {
		return nil, err
	}
	return body, nil
}

func RevokeCredential(credentialID string) error {
	_, err := makeRequest(fmt.Sprintf("/cloud-agent/credential-status/revoke-credential/%s", credentialID), "PATCH", H{})
	return err
}

func makeRequest(path string, method string, body H) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", config.Config.Wallet.Agent, path)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wallet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}

func getRequest(path string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s%s?t=%d", config.Config.Wallet.Agent, path, time.Now().Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", config.Config.Wallet.AgentApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	return respBody.Bytes(), nil
}
