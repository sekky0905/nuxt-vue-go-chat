package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/model"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/interface/controller"
)

// CheckAuthentication checks authentication of user who requested.
func CheckAuthentication() gin.HandlerFunc {
	return func(g *gin.Context) {
		id, err := g.Cookie(model.SessionIDAtCookie)
		if err != nil {
			controller.ResponseAndLogError(g, &model.AuthenticationErr{})
			return
		}

		ctx := g.Request.Context()
		repo := db.NewSessionRepository()
		m := db.NewDBManager()
		session, err := repo.GetSessionByID(ctx, m, id)
		if err != nil || session == nil {
			controller.ResponseAndLogError(g, &model.AuthenticationErr{})
			return
		}
		g.Next()
	}
}
