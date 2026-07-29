package main

import (
	"net/http"

	"builder/routes"
	"builder/store"
	"builder/utils"
)

func main() {
	server := http.FileServer(http.Dir("../public"))
	http.Handle("/", server)

	err := store.InitState()
	if err != nil {
		utils.LogError("failed to init state -", err)
		return
	}

	routes.RegisterRoutes()

	utils.LogDebug("requiem root:", store.Root)
	utils.LogInfo("listening on http://localhost:" + store.PORT)

	err = http.ListenAndServe(":"+store.PORT, nil)
	if err != nil {
		utils.LogError("failed to start server -", err)
		return
	}
}
