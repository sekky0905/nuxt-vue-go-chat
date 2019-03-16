package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ResponseAndLogError returns response and log error.
func ResponseAndLogError(g *gin.Context, err error) {
	he := handleError(err)
	if he.BaseError != nil {
		log.Errorf("error message: %v, base error: %v", he.Message, he.BaseError.Error())
	} else {
		log.Errorf("error message: %v", he.Message)
	}

	log.Infof("=====HE====%#v", he)

	g.JSON(he.Status, he)
}
