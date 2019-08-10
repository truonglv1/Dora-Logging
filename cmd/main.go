package main

import "github.com/Dora-Logs/internal/dlogs"

func main() {
	server, err := dlogs.InitServerLogging("")
	if err != nil {
		panic(err)
	}
	server.ListenAndServe()
}
