package main

import (
	"log"

	"github.com/kasif-apps/backend/auth"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(auth.AuthRoute(*app))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
