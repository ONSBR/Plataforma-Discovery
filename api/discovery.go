package api

import (
	"io/ioutil"

	"github.com/labstack/echo"
)

func discovery(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	c.String(200, string(data))
	return nil
}
