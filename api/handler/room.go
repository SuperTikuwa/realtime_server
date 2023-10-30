package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// CreateRoom creates a room
func CreateRoom(c echo.Context) error {
	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := c.Request().Context()
	roomID := uuid.New().String()

	s := r.Set(ctx, roomID, 0, 0)
	if s.Err() != nil {
		return c.String(http.StatusInternalServerError, "failed to create room")
	}

	return c.String(http.StatusCreated, "room created")
}

// ListRoom lists rooms
func ListRoom(c echo.Context) error {
	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := c.Request().Context()
	defer r.Close()

	rooms, err := r.Keys(ctx, "*").Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to list room")
	}

	var res []ListRoomResponse
	for _, room := range rooms {
		n, err := r.Get(ctx, room).Result()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to list room")
		}

		players, err := strconv.Atoi(n)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to list room")
		}
		res = append(res, ListRoomResponse{
			RoomID:  room,
			Players: players,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// ListRoomResponse is a response for ListRoom
type ListRoomResponse struct {
	RoomID  string `json:"room_id"`
	Players int    `json:"players"`
}

func ValidateRoomID(c echo.Context) error {
	roomID := c.Param("room_id")
	if roomID == "" {
		return c.String(http.StatusBadRequest, "room_id is required")
	}

	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx := c.Request().Context()
	defer r.Close()

	_, err := r.Get(ctx, roomID).Result()
	if err != nil {
		return c.String(http.StatusBadRequest, "room_id is invalid")
	}

	return c.String(http.StatusOK, "ok")

}
