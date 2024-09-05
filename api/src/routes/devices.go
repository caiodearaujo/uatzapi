package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"whatsgoingon/helpers"
)

func Devices(c *gin.Context) {
	devices, _ := helpers.GetDeviceList()

	c.JSON(http.StatusOK, devices)
}
