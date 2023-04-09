package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

func AuthRoute(app pocketbase.PocketBase) echo.Route {
	return echo.Route{
		Method:  http.MethodGet,
		Path:    "/api/users/auth-redirect",
		Handler: AuthRedirectHandler(app),
	}
}

type AuthRedirectResponse struct {
	Code string `json:"code"`
}

func AuthRedirectHandler(app pocketbase.PocketBase) func(c echo.Context) error {
	return func(c echo.Context) error {
		clients := app.SubscriptionsBroker().Clients()
		db := app.Dao().DB()

		for _, client := range clients {
			for id := range client.Subscriptions() {
				channel := strings.Split(id, "/")

				if len(channel) == 2 {
					if channel[1] == c.QueryParam("state") {
						message := AuthRedirectResponse{Code: c.QueryParam("code")}
						jsonMessage, err := json.Marshal(message)

						if err != nil {
							return nil
						}

						client.Channel() <- subscriptions.Message{Name: id, Data: string(jsonMessage)}

						query := db.NewQuery("DELETE FROM auth_table WHERE state={:state}")
						query.Bind(dbx.Params{"state": c.QueryParam("state")})
						_, err = query.Execute()

						return err
					}
				}
			}
		}

		return nil
	}
}
