package apiv1

import (
	"github.com/gin-gonic/gin"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/api/v1/auth"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/api/v1/client"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/api/v1/server"
)

// ApplyRoutes apply routes to gin router
func ApplyRoutes(r *gin.RouterGroup, private bool) {
	v1 := r.Group("/v1.0")
	{
		if private {
			client.ApplyRoutes(v1)
			server.ApplyRoutes(v1)
		} else {
			auth.ApplyRoutes(v1)

		}
	}
}
