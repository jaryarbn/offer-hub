package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

const passwordHashCost = 12

var (
	ErrInvalidRegisterParams = errors.New("invalid registration parameters")
	ErrUsernameAlreadyExists = data.ErrUsernameAlreadyExists
	ErrInvalidLoginParams    = errors.New("invalid login parameters")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrAccountUnavailable    = errors.New("account is unavailable")
	ErrInvalidLogoutToken    = errors.New("invalid logout token")
)

type AuthData interface {
	UsernameExists(context.Context, string) (bool, error)
	CreateUser(context.Context, data.CreateUserRecord) error
	GetUserByUsername(context.Context, string) (data.UserRecord, error)
	SaveLatestToken(context.Context, string, string, time.Duration) error
	DeleteLatestToken(context.Context, string) error
}

type AuthService struct {
	data          AuthData
	hashPassword  func([]byte, int) ([]byte, error)
	newUserID     func() string
	jwtSecret     []byte
	jwtExpire     time.Duration
	tokenCacheTTL time.Duration
	singleSignOn  bool
	now           func() time.Time
}

func NewAuthService(authData AuthData, jwtConfig config.JWTConfig) (*AuthService, error) {
	if authData == nil {
		return nil, errors.New("auth data is nil")
	}
	if strings.TrimSpace(jwtConfig.Secret) == "" {
		return nil, errors.New("JWT secret is empty")
	}
	if jwtConfig.Expire <= 0 {
		return nil, errors.New("JWT expiration must be positive")
	}
	tokenCacheExpire := jwtConfig.TokenCacheExpire
	if tokenCacheExpire <= 0 {
		tokenCacheExpire = jwtConfig.Expire
	}
	if tokenCacheExpire < jwtConfig.Expire {
		return nil, errors.New("JWT token cache expiration must cover token expiration")
	}

	return &AuthService{
		data:          authData,
		hashPassword:  bcrypt.GenerateFromPassword,
		newUserID:     uuid.NewString,
		jwtSecret:     []byte(jwtConfig.Secret),
		jwtExpire:     time.Duration(jwtConfig.Expire) * time.Hour,
		tokenCacheTTL: time.Duration(tokenCacheExpire) * time.Hour,
		singleSignOn:  jwtConfig.Enable,
		now:           time.Now,
	}, nil
}

func (service *AuthService) Register(ctx context.Context, req model.PasswordRegisterReq) error {
	username := strings.TrimSpace(req.Username)
	if !validRegistrationInput(username, req.Password) {
		return ErrInvalidRegisterParams
	}

	exists, err := service.data.UsernameExists(ctx, username)
	if err != nil {
		return fmt.Errorf("check username: %w", err)
	}
	if exists {
		return ErrUsernameAlreadyExists
	}

	passwordHash, err := service.hashPassword([]byte(req.Password), passwordHashCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return ErrInvalidRegisterParams
	}
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	err = service.data.CreateUser(ctx, data.CreateUserRecord{
		UserID:   service.newUserID(),
		Username: username,
		Password: string(passwordHash),
		NickName: username,
	})
	if errors.Is(err, data.ErrUsernameAlreadyExists) {
		return ErrUsernameAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("persist user: %w", err)
	}
	return nil
}

type authClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (service *AuthService) Login(
	ctx context.Context,
	req model.PasswordLoginReq,
) (model.PasswordLoginData, error) {
	username := strings.TrimSpace(req.Username)
	if !validLoginInput(username, req.Password) {
		return model.PasswordLoginData{}, ErrInvalidLoginParams
	}

	user, err := service.data.GetUserByUsername(ctx, username)
	if errors.Is(err, data.ErrUserNotFound) {
		return model.PasswordLoginData{}, ErrInvalidCredentials
	}
	if err != nil {
		return model.PasswordLoginData{}, fmt.Errorf("get login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.PasswordLoginData{}, ErrInvalidCredentials
		}
		return model.PasswordLoginData{}, fmt.Errorf("compare login password: %w", err)
	}
	if user.UserStatus != data.UserStatusActive {
		return model.PasswordLoginData{}, ErrAccountUnavailable
	}

	token, err := service.signToken(user.UserID)
	if err != nil {
		return model.PasswordLoginData{}, fmt.Errorf("sign login token: %w", err)
	}
	if service.singleSignOn {
		if err := service.data.SaveLatestToken(ctx, user.UserID, token, service.tokenCacheTTL); err != nil {
			return model.PasswordLoginData{}, fmt.Errorf("cache latest login token: %w", err)
		}
	}

	return model.PasswordLoginData{
		Token: token,
		UserInfo: model.PasswordLoginUserInfo{
			UserID:     user.UserID,
			Username:   user.Username,
			NickName:   user.NickName,
			Avatar:     user.Avatar,
			Sex:        user.Sex,
			VIP:        user.VIP,
			Phone:      user.Phone,
			Email:      user.Email,
			UserStatus: user.UserStatus,
			UserType:   user.UserType,
		},
	}, nil
}

func (service *AuthService) signToken(userID string) (string, error) {
	issuedAt := service.now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(service.jwtExpire)),
		},
	})
	return token.SignedString(service.jwtSecret)
}

func (service *AuthService) Logout(ctx context.Context, tokenString string) error {
	if strings.TrimSpace(tokenString) == "" {
		return ErrInvalidLogoutToken
	}
	now := service.now()

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&authClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidLogoutToken
			}
			return service.jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return ErrInvalidLogoutToken
	}

	claims, ok := parsedToken.Claims.(*authClaims)
	if !ok || strings.TrimSpace(claims.UserID) == "" || claims.ExpiresAt == nil {
		return ErrInvalidLogoutToken
	}

	if claims.ExpiresAt.Time.Sub(now) <= 0 {
		return ErrInvalidLogoutToken
	}
	if service.singleSignOn {
		if err := service.data.DeleteLatestToken(ctx, claims.UserID); err != nil {
			return fmt.Errorf("delete latest login token: %w", err)
		}
	}
	return nil
}

func validRegistrationInput(username, password string) bool {
	return username != "" &&
		utf8.RuneCountInString(username) <= 50 &&
		utf8.RuneCountInString(password) >= 6
}

func validLoginInput(username, password string) bool {
	return username != "" &&
		utf8.RuneCountInString(username) <= 50 &&
		utf8.RuneCountInString(password) >= 6
}
