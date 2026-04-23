package routes

import (
	"crypto/ecdh"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"builder/store"
	"builder/utils"
)

type keyBody struct {
	Key string `json:"key"`
}

func getE2EKey() {
	http.HandleFunc("/api/get-key", func(write http.ResponseWriter, req *http.Request) {
		var body keyBody

		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode body - %v", err), http.StatusInternalServerError)
			return
		}

		decodedClientPublicKey, err := base64.StdEncoding.DecodeString(body.Key)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode key - %v", err), http.StatusInternalServerError)
			return
		}

		clientPublicKey, err := ecdh.P256().NewPublicKey(decodedClientPublicKey)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to resolve client key - %v", err), http.StatusInternalServerError)
			return
		}

		privateKey, err := utils.GenerateKey()
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to generate key pair - %v", err), http.StatusInternalServerError)
			return
		}

		store.SharedSecret, err = privateKey.ECDH(clientPublicKey)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to generate shared secret - %v", err), http.StatusInternalServerError)
			return
		}

		publicKey := privateKey.PublicKey()
		encodedPublicKey := base64.StdEncoding.EncodeToString(publicKey.Bytes())

		write.Header().Set("Content-Type", "text/plain")
		write.Write([]byte(encodedPublicKey))
	})
}
