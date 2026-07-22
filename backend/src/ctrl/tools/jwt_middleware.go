package tools

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
)

const (
	unauthorizedMessage       = "未认证"
	accountUnavailableMessage = "账号不可用"
	invalidatedTokenMessage   = "Token已失效，请重新登录"
	internalErrorMessage      = "服务器内部错误"
)

type jwtAuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type userStatusGetter func(context.Context, string) (int, error)

type latestTokenChecker func(context.Context, string, string) (bool, error)

type authErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// JWTAuthMiddleware validates a Bearer token, the user status, and the latest
// SSO token before forwarding the verified user_id in the request Header.
func JWTAuthMiddleware() gin.HandlerFunc {
	var jwtConfig config.JWTConfig
	if config.Conf != nil {
		jwtConfig = config.Conf.JWT
	}
	return newJWTAuthMiddleware(jwtConfig, getUserStatus, checkTokenInRedis)
}

// SoftJWTAuthMiddleware enriches requests carrying a current JWT. Missing or
// malformed JWTs remain anonymous; a parsed but superseded SSO token is rejected.
func SoftJWTAuthMiddleware() gin.HandlerFunc {
	var jwtConfig config.JWTConfig
	if config.Conf != nil {
		jwtConfig = config.Conf.JWT
	}
	return newSoftJWTAuthMiddleware(jwtConfig, checkTokenInRedis)
}

func newJWTAuthMiddleware(
	jwtConfig config.JWTConfig,
	getStatus userStatusGetter,
	checkLatestToken latestTokenChecker,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Header.Del("user_id")
		tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			abortAuthRequest(ctx, 401, unauthorizedMessage)
			return
		}

		claims, ok := parseJWTClaims(tokenString, jwtConfig)
		if !ok {
			abortAuthRequest(ctx, 401, unauthorizedMessage)
			return
		}

		if getStatus == nil {
			log.Printf("check JWT user status: getter is not initialized")
			abortAuthRequest(ctx, 500, internalErrorMessage)
			return
		}
		userStatus, err := getStatus(ctx.Request.Context(), claims.UserID)
		if err != nil {
			log.Printf("check JWT user status: %v", err)
			abortAuthRequest(ctx, 500, internalErrorMessage)
			return
		}
		if userStatus != data.UserStatusActive {
			abortAuthRequest(ctx, 401, accountUnavailableMessage)
			return
		}

		if jwtConfig.Enable {
			if checkLatestToken == nil {
				log.Printf("check latest JWT: checker is not initialized")
				abortAuthRequest(ctx, 500, internalErrorMessage)
				return
			}
			valid, err := checkLatestToken(ctx.Request.Context(), claims.UserID, tokenString)
			if err != nil {
				log.Printf("check latest JWT: %v", err)
				abortAuthRequest(ctx, 500, internalErrorMessage)
				return
			}
			if !valid {
				abortAuthRequest(ctx, 401, invalidatedTokenMessage)
				return
			}
		}

		// Always overwrite a client-supplied user_id with the signed claim.
		ctx.Request.Header.Set("user_id", claims.UserID)
		ctx.Next()
	}
}

func newSoftJWTAuthMiddleware(
	jwtConfig config.JWTConfig,
	checkLatestToken latestTokenChecker,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// A user_id supplied by the client is never an authenticated identity.
		ctx.Request.Header.Del("user_id")

		tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.Next()
			return
		}

		claims, valid := parseJWTClaims(tokenString, jwtConfig)
		if !valid {
			ctx.Next()
			return
		}

		if jwtConfig.Enable {
			if checkLatestToken == nil {
				log.Printf("check latest soft JWT: checker is not initialized")
				abortAuthRequestWithStatus(ctx, http.StatusInternalServerError, 500, internalErrorMessage)
				return
			}
			latest, err := checkLatestToken(ctx.Request.Context(), claims.UserID, tokenString)
			if err != nil {
				log.Printf("check latest soft JWT: %v", err)
				abortAuthRequestWithStatus(ctx, http.StatusInternalServerError, 500, internalErrorMessage)
				return
			}
			if !latest {
				abortAuthRequestWithStatus(ctx, http.StatusUnauthorized, 401, invalidatedTokenMessage)
				return
			}
		}

		ctx.Request.Header.Set("user_id", claims.UserID)
		ctx.Next()
	}
}

func getUserStatus(ctx context.Context, userID string) (int, error) {
	initializedData := data.GetData()
	if initializedData == nil {
		return 0, errors.New("data is not initialized")
	}

	user, err := initializedData.GetUserByID(ctx, userID)
	if errors.Is(err, data.ErrUserNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return user.UserStatus, nil
}

func checkTokenInRedis(ctx context.Context, userID, token string) (bool, error) {
	return data.CheckLatestToken(ctx, userID, token)
}

func parseJWTClaims(tokenString string, jwtConfig config.JWTConfig) (*jwtAuthClaims, bool) {
	if strings.TrimSpace(jwtConfig.Secret) == "" {
		return nil, false
	}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&jwtAuthClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtConfig.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return nil, false
	}

	claims, ok := parsedToken.Claims.(*jwtAuthClaims)
	if !ok || strings.TrimSpace(claims.UserID) == "" {
		return nil, false
	}
	return claims, true
}

func bearerToken(authorization string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", false
	}

	token := strings.TrimPrefix(authorization, prefix)
	if token == "" || strings.TrimSpace(token) != token || strings.ContainsAny(token, " \t\r\n") {
		return "", false
	}
	return token, true
}

func abortAuthRequest(ctx *gin.Context, code int, message string) {
	abortAuthRequestWithStatus(ctx, http.StatusOK, code, message)
}

func abortAuthRequestWithStatus(ctx *gin.Context, status, code int, message string) {
	ctx.AbortWithStatusJSON(status, authErrorResponse{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}
