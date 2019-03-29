package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/logger"
	"go.uber.org/zap"
)

// ResponseAndLogError returns response and log error.
func ResponseAndLogError(g *gin.Context, err error) {
	he := handleError(err)

	errMsgField := zap.String("error message", he.Message)

	if he.BaseError != nil {
		logger.Logger.Error("", errMsgField, zap.String("base error", he.BaseError.Error()))
	} else {
		logger.Logger.Error("", errMsgField)
	}

	g.JSON(he.Status, he)
}
