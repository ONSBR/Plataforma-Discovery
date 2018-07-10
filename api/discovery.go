package api

import (
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

func discovery(c echo.Context) error {
	systemID := c.QueryParam("systemID")
	instanceID := c.QueryParam("instanceID")
	fmt.Println(systemID, instanceID)
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	fmt.Println(data)
	c.JSON(200, []string{})
	return nil
}
