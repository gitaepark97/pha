package controller

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corpix/uarand"
	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/controller/response"
	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/service"
	mockservice "github.com/gitaepark/pha/service/mock"
	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/jwt"
	"github.com/gitaepark/pha/util/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	userAgent string
	clientIp  string
)

func init() {
	userAgent = uarand.GetRandom()
	buf := make([]byte, 4)
	ip := rand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)
	clientIp = net.IP(buf).To4().String()
}

func TestResiger(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.CustomErr{}

				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
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
			name: "휴대폰번호 미입력",
			body: gin.H{
				"password": util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("phone_number")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 휴대폰번호 입력",
			body: gin.H{
				"phone_number": util.CreateRandomInt32(1, 100),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("phone_number", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "01000000000 양식이 아닌 휴대폰번호 입력",
			body: gin.H{
				"phone_number": util.CreateRandomString(10),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrPhoneNumber("phone_number")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "비밀번호 미입력",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("password")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 비밀번호 타입 입력",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomInt32(1, 10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("password", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
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

			url := "/api/auth/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestLogin(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mockService *mockservice.MockService) service.CustomErr
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				accessToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
				refreshToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				err := service.CustomErr{}

				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.LoginResponseBody{
						AccessToken:  accessToken,
						RefreshToken: refreshToken,
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
			name: "휴대폰번호 미입력",
			body: gin.H{
				"password": util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("phone_number")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 휴대폰번호 입력",
			body: gin.H{
				"phone_number": util.CreateRandomInt32(1, 100),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("phone_number", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "01000000000 양식이 아닌 휴대폰번호 입력",
			body: gin.H{
				"phone_number": util.CreateRandomString(10),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrPhoneNumber("phone_number")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "비밀번호 미입력",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("password")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 비밀번호 타입 입력",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomInt32(1, 10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(0)

				return service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("password", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			body: gin.H{
				"phone_number": util.CreateRandomPhoneNumber(),
				"password":     util.CreateRandomString(10),
			},
			buildStubs: func(mockService *mockservice.MockService) service.CustomErr {
				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					Login(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.LoginResponseBody{}, err)

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

			url := "/api/auth/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}

func TestRenewAccessToken(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)

	testCases := []struct {
		name          string
		body          func(refreshToken string) gin.H
		buildStubs    func(mockService *mockservice.MockService) (string, service.CustomErr)
		checkResponse func(recorder *httptest.ResponseRecorder, errService service.CustomErr)
	}{
		{
			name: "성공",
			body: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			buildStubs: func(mockService *mockservice.MockService) (string, service.CustomErr) {
				accessToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.AccessTokenDuration)
				refreshToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				err := service.CustomErr{}

				mockService.EXPECT().
					RenewAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.LoginResponseBody{
						AccessToken:  accessToken,
						RefreshToken: refreshToken,
					}, err)

				return refreshToken, err
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
			name: "refresh 토큰 미입력",
			body: func(refreshToken string) gin.H {
				return gin.H{}
			},
			buildStubs: func(mockService *mockservice.MockService) (string, service.CustomErr) {
				refreshToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockService.EXPECT().
					RenewAccessToken(gomock.Any(), gomock.Any()).
					Times(0)

				return refreshToken, service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(validator.ErrRequired("refresh_token")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "string 타입이 아닌 refresh 토큰 입력",
			body: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": util.CreateRandomInt32(1, 10),
				}
			},
			buildStubs: func(mockService *mockservice.MockService) (string, service.CustomErr) {
				refreshToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockService.EXPECT().
					RenewAccessToken(gomock.Any(), gomock.Any()).
					Times(0)

				return refreshToken, service.CustomErr{}
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, errService service.CustomErr) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				responseBody := getResponseBody(t, recorder.Body)
				require.Equal(t, responseBody.Meta.Code, http.StatusBadRequest)
				require.Equal(t, responseBody.Meta.Message, service.NewErrBadRequest(response.ErrType("refresh_token", "string")).Err.Error())
				require.Nil(t, responseBody.Data)
			},
		},
		{
			name: "Internal Service Error",
			body: func(refreshToken string) gin.H {
				return gin.H{
					"refresh_token": refreshToken,
				}
			},
			buildStubs: func(mockService *mockservice.MockService) (string, service.CustomErr) {
				refreshToken, _, _ := jwt.CreateToken(userID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				err := service.NewErrInternalServer(sql.ErrConnDone)

				mockService.EXPECT().
					RenewAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(dto.RenewAccessTokenResponse{}, err)

				return refreshToken, err
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

			refreshToken, errService := tc.buildStubs(service)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body(refreshToken))
			require.NoError(t, err)

			url := "/api/auth/token"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			controller.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, errService)
		})
	}
}
