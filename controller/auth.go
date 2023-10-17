package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/controller/response"
	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/service"
)

func (controller *Controller) setAuthRouter() {
	authRouter := controller.router.Group("/api/auth")

	// 회원가입 api
	authRouter.POST("/register", func(ctx *gin.Context) {
		var reqBody dto.RegisterRequestBody
		// req body dto 검증
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqBody, "json")
			return
		}

		params := service.RegisterParams(reqBody)

		// 회원가입
		cErr := controller.service.Register(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, nil)
	})

	// 로그인 api
	authRouter.POST("/login", func(ctx *gin.Context) {
		var reqBody dto.LoginRequestBody
		// req body dto 검증
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqBody, "json")
			return
		}

		params := service.LoginParams{
			LoginRequestBody: reqBody,
			UserAgent:        ctx.Request.UserAgent(),
			ClientIp:         ctx.ClientIP(),
		}

		// 로그인
		result, cErr := controller.service.Login(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, result)
	})

	// access 토큰 재발급 api
	authRouter.POST("/token", func(ctx *gin.Context) {
		var reqBody dto.RenewAccessTokenRequestBody
		// req body dto 검증
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqBody, "json")
			return
		}

		params := service.RenewAccessTokenParams{
			RenewAccessTokenRequestBody: reqBody,
			UserAgent:                   ctx.Request.UserAgent(),
			ClientIp:                    ctx.ClientIP(),
		}

		// access 토큰 재발급
		result, cErr := controller.service.RenewAccessToken(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, result)
	})
}
