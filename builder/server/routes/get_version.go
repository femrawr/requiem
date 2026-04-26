package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"builder/store"
)

type version struct {
	Server int `json:"server"`
	Public int `json:"public"`
}

func getVersion() {
	http.HandleFunc("/api/get-version", func(write http.ResponseWriter, req *http.Request) {
		versionFile := filepath.Join(store.Root, "builder", "version.json")

		if store.DEBUG {
			fmt.Println("builder versions file - " + versionFile)
		}

		data, err := os.ReadFile(versionFile)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to read version file - %v", err), http.StatusInternalServerError)
			return
		}

		var versions version

		err = json.Unmarshal(data, &versions)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to decode versions - %v", err), http.StatusInternalServerError)
			return
		}

		write.Header().Set("Content-Type", "application/json")
		json.NewEncoder(write).Encode(versions)
	})
}
