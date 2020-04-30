package main

import (
	"fmt"
	"github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/api"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/auth"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/core"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/version"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	cacheDb = cache.New(60*time.Minute, 10*time.Minute)
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Infof("Starting Wg Gen Web version: %s", version.Version)

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
	config.AddAllowHeaders("Authorization")
	app.Use(cors.New(config))

	// protection middleware
	app.Use(helmet.Default())

	// add cache storage to gin app
	app.Use(func(ctx *gin.Context) {
		ctx.Set("cache", cacheDb)
		ctx.Next()
	})

	// serve static files
	app.Use(static.Serve("/", static.LocalFile("./ui/dist", false)))

	// setup Oauth2 client
	oauth2Client, err := auth.GetAuthProvider()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to setup Oauth2")
	}

	app.Use(func(ctx *gin.Context) {
		ctx.Set("oauth2Client", oauth2Client)
		ctx.Next()
	})

	// apply api routes public
	api.ApplyRoutes(app, false)

	// simple middleware to check auth
	app.Use(func(c *gin.Context) {
		cacheDb := c.MustGet("cache").(*cache.Cache)

		token := c.Request.Header.Get(util.AuthTokenHeaderName)

		oauth2Token, exists := cacheDb.Get(token)
		if exists && oauth2Token.(*oauth2.Token).AccessToken == token {
			// will be accessible in auth endpoints
			c.Set("oauth2Token", oauth2Token)
			c.Next()
			return
		}

		// avoid 401 page for refresh after logout
		if !strings.Contains(c.Request.URL.Path, "/api/") {
			c.Redirect(301, "/index.html")
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
		return
	})

	// apply api router private
	api.ApplyRoutes(app, true)

	err = app.Run(fmt.Sprintf("%s:%s", os.Getenv("SERVER"), os.Getenv("PORT")))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to start server")
	}
}
