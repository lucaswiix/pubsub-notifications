package cors

import "github.com/gin-gonic/gin"

type Cors struct {
}

func NewCors() *Cors {
	return &Cors{}
}

// CORS will handle the CORS middleware
func (c *Cors) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With, Authorization, Accept, Origin, Cache-Control, x-api-key, x-tenant-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
