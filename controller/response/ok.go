package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewOkResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"code":    http.StatusOK,
			"message": "ok",
		},
		"data": data,
	})
}
