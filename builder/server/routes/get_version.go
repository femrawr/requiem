package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"builder/store"
)

func getVersion() {
	http.HandleFunc("/api/get-version", func(write http.ResponseWriter, req *http.Request) {
		versions := make(map[string]string)
		versions["server"] = strconv.Itoa(store.VERSION_SERVER)
		versions["public"] = strconv.Itoa(store.VERSION_PUBLIC)

		write.Header().Set("Content-Type", "application/json")
		json.NewEncoder(write).Encode(versions)
	})
}
