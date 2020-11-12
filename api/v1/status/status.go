package status

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/core"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/status")
	{
		g.GET("/enabled", readEnabled)
		g.GET("/interface", readInterfaceStatus)
		g.GET("/clients", readClientStatus)
	}
}

func readEnabled(c *gin.Context) {
	c.JSON(http.StatusOK, os.Getenv("WG_STATS_API") != "")
}

func readInterfaceStatus(c *gin.Context) {
	status, err := core.ReadInterfaceStatus()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read interface status")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, status)
}

func readClientStatus(c *gin.Context) {
	status, err := core.ReadClientStatus()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client status")
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, status)
}
