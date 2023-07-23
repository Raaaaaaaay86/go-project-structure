package main

import (
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
	"log"
)

func main() {
	err := Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	app, err := App()
	if err != nil {
		log.Fatal(err)
	}

	if err = app.HttpServer.Run(); err != nil {
		return err
	}

	return nil
}
