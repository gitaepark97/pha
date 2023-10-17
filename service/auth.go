package service

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gitaepark/pha/dto"
	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/util/bcrypt"
	"github.com/gitaepark/pha/util/jwt"
	"github.com/go-sql-driver/mysql"
)

type RegisterParams = dto.RegisterRequestBody

func (service *service) Register(ctx context.Context, params RegisterParams) (cErr CustomErr) {
	// 비밀번호 암호화
	hashedPassword, err := bcrypt.HashPassword(params.Password)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	arg := repository.CreateUserParams{
		PhoneNumber:    params.PhoneNumber,
		HashedPassword: hashedPassword,
	}

	// 회원 생성
	err = service.repository.CreateUser(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			// 휴대폰 번호가 중복된 경우
			case repository.DB_DUPLICATE_ERROR:
				switch true {
				case strings.Contains(mysqlErr.Message, "phone_number"):
					cErr = errDuplicatePhoneNumber
					return
				}
			}
		}

		cErr = NewErrInternalServer(err)
		return
	}

	return
}

type LoginParams struct {
	dto.LoginRequestBody
	UserAgent string
	ClientIp  string
}

// 로그인 로직
func (service *service) Login(ctx context.Context, params LoginParams) (result dto.LoginResponseBody, cErr CustomErr) {
	// 회원 검색
	user, err := service.repository.GetUser(ctx, params.PhoneNumber)
	if err != nil {
		// 해당 휴대폰 번호의 회원이 없는 경우
		if err == sql.ErrNoRows {
			cErr = errNotFoundUser
			return
		}
		cErr = NewErrInternalServer(err)
		return
	}

	// 비밀번호 검증 로직
	err = bcrypt.CheckPassword(params.Password, user.HashedPassword)
	if err != nil {
		cErr = errWrongPassword
		return
	}

	// access 토큰 생성
	accessToken, _, err := jwt.CreateToken(user.ID, service.config.JWTSecret, service.config.AccessTokenDuration)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	// refresh 토큰 생성
	refreshToken, refreshPayload, err := jwt.CreateToken(user.ID, service.config.JWTSecret, service.config.RefreshTokenDuration)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	arg := repository.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    params.UserAgent,
		ClientIp:     params.ClientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	// 세션 저장
	err = service.repository.CreateSession(ctx, arg)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	result = dto.LoginResponseBody{AccessToken: accessToken, RefreshToken: refreshToken}
	return
}

type RenewAccessTokenParams struct {
	dto.RenewAccessTokenRequestBody
	UserAgent string
	ClientIp  string
}

// access 토큰 재발급 로직
func (service *service) RenewAccessToken(ctx context.Context, params RenewAccessTokenParams) (result dto.RenewAccessTokenResponse, cErr CustomErr) {
	// 토큰 검증
	refreshPayload, err := jwt.VerifyToken(params.RefreshToken, service.config.JWTSecret)
	if err != nil {
		cErr = NewErrBadRequest(err)
		return
	}

	// 세션 검색
	session, err := service.repository.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		// 해당 id의 세션이 없는 경우
		if err == sql.ErrNoRows {
			cErr = errNotFoundSession
			return
		}

		cErr = NewErrInternalServer(err)
		return
	}

	// 세션 검증
	// 세션이 막힌 경우
	if session.IsBlocked {
		cErr = errBlockedSession
		return
	}
	// 세션 회원이 아닌 경우
	if session.UserID != refreshPayload.UserID {
		cErr = errIncorrectSessionUser
		return
	}
	// refresh 토큰이 일치하지 않는 경우
	if session.RefreshToken != params.RefreshToken {
		cErr = errMismatchedSessionToken
		return
	}
	// 세션이 만료된 경우
	if time.Now().After(session.ExpiredAt) {
		cErr = errExpiredSession
		return
	}

	// access 토큰 생성
	accessToken, _, err := jwt.CreateToken(refreshPayload.UserID, service.config.JWTSecret, service.config.AccessTokenDuration)
	if err != nil {
		cErr = NewErrInternalServer(err)
		return
	}

	result = dto.RenewAccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: params.RefreshToken,
	}

	return
}
