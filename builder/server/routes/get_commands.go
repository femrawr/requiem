package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"builder/store"
)

func getCommands() {
	http.HandleFunc("/api/get-commands", func(write http.ResponseWriter, req *http.Request) {
		commands := filepath.Join(store.Main, "bot", "commands")

		items, err := os.ReadDir(commands)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to read dir - %s", err), http.StatusInternalServerError)
			return
		}

		type Command struct {
			Name string `json:"name"`
			Info string `json:"info"`
		}

		var result []Command
		nameRegex := regexp.MustCompile(`func \(.*?\) Name\(\) string \{\s*return "(.+?)"`)
		infoRegex := regexp.MustCompile(`func \(.*?\) Info\(\) string \{\s*return "(.+?)"`)

		for _, item := range items {
			if item.IsDir() {
				continue
			}

			data, err := os.ReadFile(filepath.Join(commands, item.Name()))
			if err != nil {
				continue
			}

			content := string(data)

			nameMatch := nameRegex.FindStringSubmatch(content)
			infoMatch := infoRegex.FindStringSubmatch(content)
			if len(nameMatch) < 2 || len(infoMatch) < 2 {
				continue
			}

			result = append(result, Command{
				Name: nameMatch[1],
				Info: infoMatch[1],
			})
		}

		write.Header().Set("Content-Type", "application/json")
		json.NewEncoder(write).Encode(result)
	})
}
