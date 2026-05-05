package main

import (
	"github.com/notOliveira/onde-tem/internal/infra/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.Close()

	app.Router.Run(app.Config.ServerPort)
}
