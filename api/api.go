package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/repository"
	"net/http"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.Engine) {
	client := r.Group("/api/v1.0/client")
	{

		client.POST("", createClient)
		client.GET("/:id", readClient)
		client.PATCH("/:id", updateClient)
		client.DELETE("/:id", deleteClient)
		client.GET("", readClients)
		client.GET("/:id/config", configClient)
	}

	server := r.Group("/api/v1.0/server")
	{
		server.GET("", readServer)
		server.PATCH("", updateServer)
	}
}


func createClient(c *gin.Context) {
	var data model.Client

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	client, err := repository.CreateClient(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func readClient(c *gin.Context) {
	id := c.Param("id")

	client, err := repository.ReadClient(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateClient(c *gin.Context) {
	var data model.Client
	id := c.Param("id")

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	client, err := repository.UpdateClient(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteClient(c *gin.Context) {
	id := c.Param("id")

	err := repository.DeleteClient(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to remove client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func readClients(c *gin.Context) {
	clients, err := repository.ReadClients()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to list clients")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, clients)
}

func configClient(c *gin.Context) {
	configData, err := repository.ReadClientConfig(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client config")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	formatQr := c.DefaultQuery("qrcode", "false")
	if formatQr == "false" {
		// return config as txt file
		c.Header("Content-Disposition", "attachment; filename=wg0.conf")
		c.Data(http.StatusOK, "application/config", configData)
		return
	} else {
		// return config as png qrcode
		png, err := qrcode.Encode(string(configData), qrcode.Medium, 220)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to create qrcode")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "image/png", png)
		return
	}
}

func readServer(c *gin.Context) {
	client, err := repository.ReadServer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateServer(c *gin.Context) {
	var data model.Server

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	client, err := repository.UpdateServer(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}