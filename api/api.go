package api

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func InitAPI() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = errorHandler
	group := e.Group("/v1.0.0")
	group.POST("/discovery", discovery)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func errorHandler(err error, c echo.Context) {
	if errJ := c.JSON(400, map[string]string{"status": "400", "message": err.Error()}); errJ != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}
