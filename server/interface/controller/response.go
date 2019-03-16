package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ResponseAndLogError returns response and log error.
func ResponseAndLogError(g *gin.Context, err error) {
	he := handleError(err)
	log.Errorf("error message: %v, base error: %v", he.Message, he.BaseError.Error())

	g.JSON(he.Status, he)
}
