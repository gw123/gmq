package webMiddlewares

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/db_models"
	"github.com/labstack/echo"
	"strconv"
	"time"
)

func NewPingMiddleware(app interfaces.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			clientId := c.QueryParam("clientId")
			payload := c.QueryParam("payload")
			sendAt := c.QueryParam("sendAt")

			sendAtTIme, _ := time.Parse("2019-08-24", sendAt)
			if sendAtTIme.IsZero() {
				sendAtTIme = time.Now()
			}
			l := stop.Sub(start)
			if l != 0 {
				l = l / time.Millisecond
			}

			byteIn, _ := strconv.Atoi(req.Header.Get(echo.HeaderContentLength))
			pingLog := &db_models.PingLog{

				Ip:           c.RealIP(),
				ClientSendAt: sendAtTIme,
				CreatedAt:    time.Now(),
				Payload:      payload,
				ClientId:     clientId,
				Latency:      uint(l),
				BytesIn:      uint(byteIn),
				BytesOut:     uint(res.Size),
			}
			db, err := app.GetDefaultDb()
			if err != nil {
			} else {
				if err = db.Save(pingLog).Error; err != nil {
				}
			}
			return nil
		}
	}
}
