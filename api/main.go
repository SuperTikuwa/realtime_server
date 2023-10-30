package main

import (
	"github.com/labstack/echo/v4"
	"github.com/supertikuwa/realtime_server/api/handler"
)

func main() {
	e := echo.New()
	g := e.Group("/api")

	g.GET("/hc", func(c echo.Context) error {
		return c.String(200, "healthy")
	})

	g.GET("/room", handler.ListRoom)
	g.POST("/room", handler.CreateRoom)
	g.GET("/room/:room_id", handler.ValidateRoomID)

	e.Logger.Fatal(e.Start(":80"))
}
