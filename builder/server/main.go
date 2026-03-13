package main

import (
	"fmt"
	"net/http"

	"builder/routes"
	"builder/store"
)

func main() {
	server := http.FileServer(http.Dir("../public"))
	http.Handle("/", server)

	store.InitState()
	routes.RegisterRoutes()

	if store.DEBUG {
		fmt.Println("requiem root - " + store.Root)
	}

	fmt.Println("listening on http://localhost:" + store.PORT)

	err := http.ListenAndServe(":"+store.PORT, nil)
	if err != nil {
		fmt.Printf("failed to start server: %s", err)
		return
	}
}
