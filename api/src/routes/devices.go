package routes

import (
	"net/http"
	"whatsgoingon/helpers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Devices(c *gin.Context) {
	devices, _ := helpers.GetDeviceList()

	c.JSON(http.StatusOK, devices)
}
