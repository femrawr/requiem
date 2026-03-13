package routes

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"builder/store"
)

func startBuild() {
	http.HandleFunc("/api/start-build", func(write http.ResponseWriter, req *http.Request) {
		buildDir := filepath.Join(store.Root, ".builds")
		buildName := fmt.Sprintf("%s-%d.exe", store.Tag, time.Now().UnixNano())
		buildPath := filepath.Join(buildDir, buildName)

		cmd := exec.Command(
			"go", "build",
			"-trimpath",
			"-buildvcs=false",
			"-ldflags=-s -w -H windowsgui -buildid=",
			"-o", buildPath,
		)

		cmd.Dir = store.Main
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		err := cmd.Run()
		if err != nil {
			http.Error(write, fmt.Sprintf("failed to build - %s", err), http.StatusInternalServerError)
			return
		}

		write.WriteHeader(http.StatusOK)
	})
}
