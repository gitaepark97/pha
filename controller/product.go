package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/controller/response"
	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/middleware"
	"github.com/gitaepark/pha/service"
	"github.com/gitaepark/pha/util/jwt"
)

func (controller *Controller) setProductRouter() {
	// authorization
	productRoutes := controller.router.Group("/api/products").Use(middleware.AuthMiddleware(controller.config))

	// 상품 등록 api
	productRoutes.POST("/", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*jwt.Payload)

		var reqBody dto.CreateProductRequestBody
		// req body dto 검증
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqBody, "json")
			return
		}

		params := service.CreateProductParams{
			UserID:                   authPayload.UserID,
			CreateProductRequestBody: reqBody,
		}

		// 상품 등록
		cErr := controller.service.CreateProduct(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, nil)
	})

	// 상품 목록 조회 api
	productRoutes.GET("/", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*jwt.Payload)

		var reqQuery dto.GetProductListRequestQuery
		// req query dto 검증
		if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqQuery, "form")
			return
		}

		params := service.GetProductListParams{
			UserID:                     authPayload.UserID,
			GetProductListRequestQuery: reqQuery,
		}

		// 상품 목록 조회
		result, cErr := controller.service.GetProductList(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, result)
	})

	// 상품 상세 조회 api
	productRoutes.GET("/:id", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*jwt.Payload)

		var reqPath dto.GetProductRequestPath
		// req path dto 검증
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqPath, "uri")
			return
		}

		params := service.GetProductParams{
			UserID:                authPayload.UserID,
			GetProductRequestPath: reqPath,
		}

		// 상품 상세 조회
		result, cErr := controller.service.GetProduct(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, result)
	})

	// 상품 수정 api
	productRoutes.PATCH("/:id", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*jwt.Payload)

		var reqPath dto.UpdateProductRequestPath
		// req path dto 검증
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqPath, "uri")
			return
		}
		var reqBody dto.UpdateProductRequestBody
		// req body dto 검증
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqBody, "json")
			return
		}

		params := service.UpdateProductParams{
			UserID:                   authPayload.UserID,
			UpdateProductRequestPath: reqPath,
			UpdateProductRequestBody: reqBody,
		}

		cErr := controller.service.UpdateProduct(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, nil)
	})

	productRoutes.DELETE("/:id", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*jwt.Payload)

		var reqPath dto.DeleteProductRequestPath
		// req path dto 검증
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			response.NewErrBindingResponse(ctx, err, &reqPath, "uri")
			return
		}

		params := service.DeleteProductParams{
			UserID:                   authPayload.UserID,
			DeleteProductRequestPath: reqPath,
		}

		// 상품 삭제
		cErr := controller.service.DeleteProduct(ctx, params)
		if cErr.Err != nil {
			response.NewErrResponse(ctx, cErr)
			return
		}

		response.NewOkResponse(ctx, nil)
	})
}
