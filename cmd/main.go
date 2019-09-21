package main

import "github.com/Dora-Logging/internal/dlogs"

func main() {
	server, err := dlogs.InitServerLogging("")
	if err != nil {
		panic(err)
	}
	server.ListenAndServe()
}
