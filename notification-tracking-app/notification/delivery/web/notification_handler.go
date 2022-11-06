package web

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lucaswiix/notifications-tracking-app/domain"
)

type NotificationHandler struct {
	upgrader websocket.Upgrader
	NUsecase domain.NotificationUsecase
}

func NewNotificationHandler(nu domain.NotificationUsecase) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "website/index.html")

	handler := &NotificationHandler{
		upgrader: websocket.Upgrader{},
		NUsecase: nu,
	}

	e.GET("/notifications/track/:userID", handler.TrackByUserID)

	e.Logger.Fatal(e.Start(":1323"))
}

func (p *NotificationHandler) TrackByUserID(c echo.Context) error {
	wsConn, err := p.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		_, _, err = wsConn.ReadMessage()
		if err != nil {
			cancelFunc()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wsConn.Close()
			return nil
		default:
			p, err := p.NUsecase.TrackByUserID(ctx, c.Param("userID"), "web")
			if err != nil {
				c.Logger().Error(err)
				continue
			}

			err = wsConn.WriteJSON(p)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}
