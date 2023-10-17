package service

import (
	"context"
	"database/sql"
	"encoding/binary"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/corpix/uarand"
	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/repository"
	mockrepository "github.com/gitaepark/pha/repository/mock"
	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/bcrypt"
	"github.com/gitaepark/pha/util/jwt"
	"github.com/go-sql-driver/mysql"
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

func TestRegister(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		params        RegisterParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(err CustomErr)
	}{
		{
			name: "성공",
			params: RegisterParams{
				PhoneNumber: user.PhoneNumber,
				Password:    password,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err CustomErr) {
				require.Empty(t, err)
			},
		},
		{
			name: "중복된 휴대폰번호",
			params: RegisterParams{
				PhoneNumber: user.PhoneNumber,
				Password:    password,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mysql.MySQLError{
						Number:  repository.DB_DUPLICATE_ERROR,
						Message: "phone_number",
					})
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, errDuplicatePhoneNumber)
			},
		},
		{
			name: "Internal Server Error",
			params: RegisterParams{
				PhoneNumber: user.PhoneNumber,
				Password:    password,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(err CustomErr) {
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			err := service.Register(context.Background(), tc.params)
			tc.checkResponse(err)
		})
	}
}

func TestLogin(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []struct {
		name          string
		params        LoginParams
		buildStubs    func(mockRepository *mockrepository.MockRepository)
		checkResponse func(result dto.LoginResponseBody, err CustomErr)
	}{
		{
			name: "성공",
			params: LoginParams{
				LoginRequestBody: dto.LoginRequestBody{
					PhoneNumber: user.PhoneNumber,
					Password:    password,
				},
				UserAgent: userAgent,
				ClientIp:  clientIp,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
				mockRepository.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				payload, _ := jwt.VerifyToken(result.AccessToken, testConfig.JWTSecret)
				require.Equal(t, user.ID, payload.UserID)
				require.Empty(t, err)
			},
		},
		{
			name: "회원이 없는 경우",
			params: LoginParams{
				LoginRequestBody: dto.LoginRequestBody{
					PhoneNumber: user.PhoneNumber,
					Password:    password,
				},
				UserAgent: userAgent,
				ClientIp:  clientIp,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.User{}, sql.ErrNoRows)
				mockRepository.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errNotFoundUser)
			},
		},
		{
			name: "비밀번호가 틀린 경우",
			params: LoginParams{
				LoginRequestBody: dto.LoginRequestBody{
					PhoneNumber: user.PhoneNumber,
					Password:    util.CreateRandomString(10),
				},
				UserAgent: userAgent,
				ClientIp:  clientIp,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
				mockRepository.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errWrongPassword)
			},
		},
		{
			name: "Internal Server Error",
			params: LoginParams{
				LoginRequestBody: dto.LoginRequestBody{
					PhoneNumber: user.PhoneNumber,
					Password:    password,
				},
				UserAgent: userAgent,
				ClientIp:  clientIp,
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) {
				mockRepository.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.User{}, sql.ErrConnDone)
				mockRepository.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			tc.buildStubs(repository)

			result, err := service.Login(context.Background(), tc.params)
			tc.checkResponse(result, err)
		})
	}
}

func TestRenewAccessToken(t *testing.T) {
	user, _ := createRandomUser(t)

	testCases := []struct {
		name          string
		params        func(refreshToken string) RenewAccessTokenParams
		buildStubs    func(mockRepository *mockrepository.MockRepository) string
		checkResponse func(result dto.LoginResponseBody, err CustomErr)
	}{
		{
			name: "성공",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{
						ID:           refreshPayload.ID,
						UserID:       user.ID,
						RefreshToken: refreshToken,
						UserAgent:    userAgent,
						ClientIp:     clientIp,
						IsBlocked:    false,
						ExpiredAt:    refreshPayload.ExpiredAt,
					}, nil)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				payload, _ := jwt.VerifyToken(result.AccessToken, testConfig.JWTSecret)
				require.Equal(t, user.ID, payload.UserID)
				require.Empty(t, err)
			},
		},
		{
			name: "유효하지 않은 refresh 토큰",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken := util.CreateRandomString(50)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, NewErrBadRequest(jwt.ErrInvalidToken))
			},
		},
		{
			name: "만료된 refresh 토큰",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, _, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, -time.Minute)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, NewErrBadRequest(jwt.ErrExpiredToken))
			},
		},
		{
			name: "세션이 없는 경우",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, _, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{}, sql.ErrNoRows)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errNotFoundSession)
			},
		},
		{
			name: "세션이 막힌 경우",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{
						ID:           refreshPayload.ID,
						UserID:       user.ID,
						RefreshToken: refreshToken,
						UserAgent:    userAgent,
						ClientIp:     clientIp,
						IsBlocked:    true,
						ExpiredAt:    refreshPayload.ExpiredAt,
					}, nil)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errBlockedSession)
			},
		},
		{
			name: "세션 회원이 아닌 경우",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{
						ID:           refreshPayload.ID,
						UserID:       util.CreateRandomInt64(11, 20),
						RefreshToken: refreshToken,
						UserAgent:    userAgent,
						ClientIp:     clientIp,
						IsBlocked:    false,
						ExpiredAt:    refreshPayload.ExpiredAt,
					}, nil)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errIncorrectSessionUser)
			},
		},
		{
			name: "refresh 토큰이 일치하지 않는 경우",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{
						ID:           refreshPayload.ID,
						UserID:       user.ID,
						RefreshToken: util.CreateRandomString(50),
						UserAgent:    userAgent,
						ClientIp:     clientIp,
						IsBlocked:    false,
						ExpiredAt:    refreshPayload.ExpiredAt,
					}, nil)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errMismatchedSessionToken)
			},
		},
		{
			name: "refresh 토큰이 일치하지 않는 경우",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{
						ID:           refreshPayload.ID,
						UserID:       user.ID,
						RefreshToken: refreshToken,
						UserAgent:    userAgent,
						ClientIp:     clientIp,
						IsBlocked:    false,
						ExpiredAt:    time.Now().Add(-time.Minute),
					}, nil)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, errExpiredSession)
			},
		},
		{
			name: "Internal Server Error",
			params: func(refreshToken string) RenewAccessTokenParams {
				return RenewAccessTokenParams{
					RenewAccessTokenRequestBody: dto.RenewAccessTokenRequestBody{
						RefreshToken: refreshToken,
					},
					UserAgent: userAgent,
					ClientIp:  clientIp,
				}
			},
			buildStubs: func(mockRepository *mockrepository.MockRepository) string {
				refreshToken, _, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

				mockRepository.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(repository.Session{}, sql.ErrConnDone)

				return refreshToken
			},
			checkResponse: func(result dto.LoginResponseBody, err CustomErr) {
				require.Empty(t, result)
				require.Equal(t, err, NewErrInternalServer(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockrepository.NewMockRepository(ctrl)
			service := newTestService(t, repository)

			refreshToken := tc.buildStubs(repository)

			result, err := service.RenewAccessToken(context.Background(), tc.params(refreshToken))
			tc.checkResponse(result, err)
		})
	}
}

func createRandomUser(t *testing.T) (repository.User, string) {
	password := util.CreateRandomString(10)
	hashedPassword, _ := bcrypt.HashPassword(password)

	user := repository.User{
		ID:             util.CreateRandomInt64(1, 10),
		PhoneNumber:    util.CreateRandomPhoneNumber(),
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	return user, password
}
