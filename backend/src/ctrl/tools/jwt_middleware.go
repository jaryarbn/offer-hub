package tools

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
)

const (
	unauthorizedMessage  = "未认证"
	internalErrorMessage = "服务器内部错误"
)

type jwtAuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type tokenBlacklistChecker func(context.Context, string) (bool, error)

type authErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// JWTAuthMiddleware validates a Bearer token, rejects logged-out tokens, and
// forwards the verified user_id to handlers through the request Header.
func JWTAuthMiddleware() gin.HandlerFunc {
	var jwtConfig config.JWTConfig
	if config.Conf != nil {
		jwtConfig = config.Conf.JWT
	}
	return newJWTAuthMiddleware(jwtConfig, data.IsTokenBlacklisted)
}

// SoftJWTAuthMiddleware enriches a request with user_id when a valid JWT is
// present. Missing or invalid credentials are treated as an anonymous request.
func SoftJWTAuthMiddleware() gin.HandlerFunc {
	var jwtConfig config.JWTConfig
	if config.Conf != nil {
		jwtConfig = config.Conf.JWT
	}
	return newSoftJWTAuthMiddleware(jwtConfig)
}

func newJWTAuthMiddleware(
	jwtConfig config.JWTConfig,
	isTokenBlacklisted tokenBlacklistChecker,
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

		if isTokenBlacklisted == nil {
			log.Printf("check JWT blacklist: checker is not initialized")
			abortAuthRequest(ctx, 500, internalErrorMessage)
			return
		}
		blacklisted, err := isTokenBlacklisted(ctx.Request.Context(), tokenString)
		if err != nil {
			log.Printf("check JWT blacklist: %v", err)
			abortAuthRequest(ctx, 500, internalErrorMessage)
			return
		}
		if blacklisted {
			abortAuthRequest(ctx, 401, unauthorizedMessage)
			return
		}

		// Always overwrite a client-supplied user_id with the signed claim.
		ctx.Request.Header.Set("user_id", claims.UserID)
		ctx.Next()
	}
}

func newSoftJWTAuthMiddleware(jwtConfig config.JWTConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// A user_id supplied by the client is never an authenticated identity.
		ctx.Request.Header.Del("user_id")

		tokenString, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.Next()
			return
		}

		if claims, valid := parseJWTClaims(tokenString, jwtConfig); valid {
			ctx.Request.Header.Set("user_id", claims.UserID)
		}
		ctx.Next()
	}
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
	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func abortAuthRequest(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(http.StatusOK, authErrorResponse{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}
