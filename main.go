package main

import (
	"fmt"
	"github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/api"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/core"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"os"
	"path/filepath"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infof("Starting Wg Gen Web version: %s", util.Version)

	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to load .env file")
	}

	// check directories or create it
	if !util.DirectoryExists(filepath.Join(os.Getenv("WG_CONF_DIR"))) {
		err = os.Mkdir(filepath.Join(os.Getenv("WG_CONF_DIR")), 0755)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"dir": filepath.Join(os.Getenv("WG_CONF_DIR")),
			}).Fatal("failed to create directory")
		}
	}

	// check if server.json exists otherwise create it with default values
	if !util.FileExists(filepath.Join(os.Getenv("WG_CONF_DIR"), "server.json")) {
		_, err = core.ReadServer()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("server.json doesnt not exists and can not read it")
		}
	}

	if os.Getenv("GIN_MODE") == "debug" {
		// set gin release debug
		gin.SetMode(gin.DebugMode)
	} else {
		// set gin release mode
		gin.SetMode(gin.ReleaseMode)
		// disable console color
		gin.DisableConsoleColor()
		// log level info
		log.SetLevel(log.InfoLevel)
	}

	// migrate
	err = core.MigrateInitialStructChange()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to migrate initial struct changes")
	}
	err = core.MigratePresharedKey()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to migrate preshared key struct changes")
	}

	// dump wg config file
	err = core.UpdateServerConfigWg()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to dump wg config file")
	}

	// creates a gin router with default middleware: logger and recovery (crash-free) middleware
	app := gin.Default()

	// cors middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	app.Use(cors.New(config))

	// protection middleware
	app.Use(helmet.Default())

	// no route redirect to frontend app
	app.NoRoute(func(c *gin.Context) {
		c.Redirect(301, "/index.html")
	})

	// serve static files
	app.Use(static.Serve("/", static.LocalFile("./ui/dist", false)))

	// apply api router
	api.ApplyRoutes(app)

	err = app.Run(fmt.Sprintf("%s:%s", os.Getenv("SERVER"), os.Getenv("PORT")))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to start server")
	}
}
