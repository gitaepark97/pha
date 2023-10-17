package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/controller/response"
	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/middleware"
	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/service"
	mockservice "github.com/gitaepark/pha/service/mock"
	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/jwt"
	"github.com/gitaepark/pha/util/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	product := createRandomProduct(userID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "카테고리 미입력",
			body: gin.H{
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("category")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 카테고리 입력",
			body: gin.H{
				"category":        util.CreateRandomInt32(1, 10),
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("category", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "길이가 100자 초과인 카테고리 입력",
			body: gin.H{
				"category":        util.CreateRandomString(101),
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrMax("category", "100")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "가격 미입력",
			body: gin.H{
				"category":        product.Category,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("price")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int32 타입이 아닌 가격 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           util.CreateRandomString(10),
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("price", "int32")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "원가 미입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("cost")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int32 타입이 아닌 원가 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            util.CreateRandomString(10),
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("cost", "int32")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "이름 미입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("name")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 이름 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            util.CreateRandomInt32(1, 10),
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("name", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "길이가 100자 초과인 이름 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            util.CreateRandomString(101),
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrMax("name", "100")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "설명 미입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("description")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 설명 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     util.CreateRandomInt32(1, 10),
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("description", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "바코드 미입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("barcode")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 바코드 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         util.CreateRandomInt32(1, 10),
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("barcode", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "유통기한 미입력",
			body: gin.H{
				"category":    product.Category,
				"price":       product.Price,
				"cost":        product.Cost,
				"name":        product.Name,
				"description": product.Description,
				"barcode":     product.Barcode,
				"size":        product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("expiration_date")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 유통기한 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": util.CreateRandomInt32(1, 10),
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("expiration_date", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "0000-00-00 양식이 아닌 유통기한 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": util.CreateRandomString(10),
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrDate("expiration_date")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "사이즈 미입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("size")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 사이즈 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("size", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "large와 small이 아닌 사이즈 입력",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            util.CreateRandomString(10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrProductSize("size")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, errService.Code)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, errService.Code)
				require.Equal(t, responseBody.Meta.Message, errService.Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockservice.NewMockService(ctrl)
			controller := newTestController(t, service)

			errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/products/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestGetProductList(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	product := createRandomProduct(userID)

	testCases := []struct {
		name          string
		uri           string
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			uri:  "?page=" + fmt.Sprint(util.CreateRandomInt32(1, 5)),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.GetProductListResponse{
						List: []dto.GetProductResponse{product},
					}, err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.NotEmpty(t, responseBody.Data)
			},
		},
		{
			name: "검색 성공",
			uri:  "?page=" + fmt.Sprint(util.CreateRandomInt32(1, 5)) + "&keyword=" + fmt.Sprint(product.Name),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.GetProductListResponse{
						List: []dto.GetProductResponse{product},
					}, err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.NotEmpty(t, responseBody.Data)
			},
		},
		{
			name: "페이지 미입력",
			uri:  "",
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(0)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("page")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int32 타입이 아닌 페이지 입력",
			uri:  "?page=" + util.CreateRandomString(5),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(0)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrParseString).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			uri:  "?page=" + fmt.Sprint(util.CreateRandomInt32(1, 5)),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					GetProductList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.GetProductListResponse{}, err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, errService.Code)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, errService.Code)
				require.Equal(t, responseBody.Meta.Message, errService.Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockservice.NewMockService(ctrl)
			controller := newTestController(t, service)

			errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			url := "/api/products/" + tc.uri
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestGetProduct(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	product := createRandomProduct(userID)

	testCases := []struct {
		name          string
		uri           string
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			uri:  fmt.Sprint(product.ID),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(product, err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.NotEmpty(t, responseBody.Data)
			},
		},
		{
			name: "int64 타입이 아닌 id",
			uri:  util.CreateRandomString(5),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrParseString).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			uri:  fmt.Sprint(product.ID),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.GetProductResponse{}, err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, errService.Code)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, errService.Code)
				require.Equal(t, responseBody.Meta.Message, errService.Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockservice.NewMockService(ctrl)
			controller := newTestController(t, service)

			errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			url := "/api/products/" + tc.uri
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	product := createRandomProduct(userID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			body: gin.H{
				"category": product.Category,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 카테고리 입력",
			body: gin.H{
				"category": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("category", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "길이가 100자 초과인 카테고리 입력",
			body: gin.H{
				"category": util.CreateRandomString(101),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrMax("category", "100")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int32 타입이 아닌 가격 입력",
			body: gin.H{
				"price": util.CreateRandomString(10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("price", "int32")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int32 타입이 아닌 원가 입력",
			body: gin.H{
				"cost": util.CreateRandomString(10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("cost", "int32")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 이름 입력",
			body: gin.H{
				"name": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("name", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "길이가 100자 초과인 이름 입력",
			body: gin.H{
				"name": util.CreateRandomString(101),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrMax("name", "100")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 설명 입력",
			body: gin.H{
				"description": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("description", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 바코드 입력",
			body: gin.H{
				"barcode": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("barcode", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 유통기한 입력",
			body: gin.H{
				"expiration_date": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("expiration_date", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "0000-00-00 양식이 아닌 유통기한 입력",
			body: gin.H{
				"expiration_date": util.CreateRandomString(10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrDate("expiration_date")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 사이즈 입력",
			body: gin.H{
				"size": util.CreateRandomInt32(1, 10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("size", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "large와 small이 아닌 사이즈 입력",
			body: gin.H{
				"size": util.CreateRandomString(10),
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrProductSize("size")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			body: gin.H{
				"category":        product.Category,
				"price":           product.Price,
				"cost":            product.Cost,
				"name":            product.Name,
				"description":     product.Description,
				"barcode":         product.Barcode,
				"expiration_date": product.ExpirationDate,
				"size":            product.Size,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, errService.Code)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, errService.Code)
				require.Equal(t, responseBody.Meta.Message, errService.Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockservice.NewMockService(ctrl)
			controller := newTestController(t, service)

			errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/products/" + fmt.Sprint(product.ID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	product := createRandomProduct(userID)

	testCases := []struct {
		name          string
		uri           string
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			uri:  fmt.Sprint(product.ID),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusOK)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusOK)
				require.Equal(t, responseBody.Meta.Message, "ok")
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "int64 타입이 아닌 id",
			uri:  util.CreateRandomString(5),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(0)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrParseString).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			uri:  fmt.Sprint(product.ID),
			setupAuth: func(t *testing.T, request *http.Request) {
				AddAuthorization(t, request, middleware.AuthorizationTypeBearer, userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(err)

				return err
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, errService.Code)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, errService.Code)
				require.Equal(t, responseBody.Meta.Message, errService.Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mockservice.NewMockService(ctrl)
			controller := newTestController(t, service)

			errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			url := "/api/products/" + tc.uri
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func AddAuthorization(t *testing.T, request *http.Request, authorizationType string, userID int64, secret string, duration time.Duration) {
	token, payload, err := jwt.CreateToken(userID, secret, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(middleware.AuthorizationHeaderKey, authorizationHeader)
}

func createRandomProduct(userID int64) dto.GetProductResponse {
	product := repository.Product{
		ID:             util.CreateRandomInt64(1, 10),
		UserID:         userID,
		Category:       util.CreateRandomString(15),
		Price:          util.CreateRandomInt32(1000, 10000),
		Cost:           util.CreateRandomInt32(1000, 10000),
		Name:           util.CreateRandomString(10),
		Description:    util.CreateRandomString(50),
		Barcode:        util.CreateRandomString(12),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Size:           repository.ProductSize(util.CreateRandomProductSize()),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return dto.NewGetProductResponse(product)
}
