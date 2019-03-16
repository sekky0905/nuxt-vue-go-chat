package router

import (
	"github.com/gin-gonic/gin"
)

// G is the gin engine.
var G *gin.Engine

func init() {
	g := gin.New()
	G = g
}
