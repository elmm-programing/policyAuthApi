package main

import (
	"fmt"
	"policyAuth/internal/server"
)

func main() {
	srv := server.NewServer()

	err := srv.Start()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

