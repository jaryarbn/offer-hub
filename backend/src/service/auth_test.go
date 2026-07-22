package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"offer-hub/backend/src/config"
	"offer-hub/backend/src/data"
	"offer-hub/backend/src/model"
)

type authDataStub struct {
	exists          bool
	existsErr       error
	createErr       error
	checkedUsername string
	record          data.CreateUserRecord
	createCalls     int
	user            data.UserRecord
	userErr         error
	loginUsername   string
	blacklistToken  string
	blacklistTTL    time.Duration
	blacklistErr    error
}

func (stub *authDataStub) UsernameExists(_ context.Context, username string) (bool, error) {
	stub.checkedUsername = username
	return stub.exists, stub.existsErr
}

func (stub *authDataStub) CreateUser(_ context.Context, record data.CreateUserRecord) error {
	stub.createCalls++
	stub.record = record
	return stub.createErr
}

func (stub *authDataStub) GetUserByUsername(_ context.Context, username string) (data.UserRecord, error) {
	stub.loginUsername = username
	return stub.user, stub.userErr
}

func (stub *authDataStub) AddTokenToBlacklist(_ context.Context, token string, ttl time.Duration) error {
	stub.blacklistToken = token
	stub.blacklistTTL = ttl
	return stub.blacklistErr
}

func newAuthServiceForTest(t *testing.T, stub AuthData) *AuthService {
	t.Helper()
	authService, err := NewAuthService(stub, config.JWTConfig{
		Secret:           "test-jwt-secret",
		Expire:           24,
		TokenCacheExpire: 48,
	})
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	return authService
}

func TestAuthServiceRegisterCreatesUser(t *testing.T) {
	stub := &authDataStub{}
	authService := newAuthServiceForTest(t, stub)
	var gotCost int
	authService.hashPassword = func(password []byte, cost int) ([]byte, error) {
		if string(password) != "123456" {
			t.Fatalf("password passed to hasher = %q", password)
		}
		gotCost = cost
		return []byte("bcrypt-hash"), nil
	}
	authService.newUserID = func() string { return "user-uuid" }

	err := authService.Register(context.Background(), model.PasswordRegisterReq{
		Username: "  testuser  ",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if stub.checkedUsername != "testuser" {
		t.Fatalf("checked username = %q, want testuser", stub.checkedUsername)
	}
	if gotCost != passwordHashCost {
		t.Fatalf("bcrypt cost = %d, want %d", gotCost, passwordHashCost)
	}
	want := data.CreateUserRecord{
		UserID:   "user-uuid",
		Username: "testuser",
		Password: "bcrypt-hash",
		NickName: "testuser",
	}
	if stub.record != want {
		t.Fatalf("created user = %#v, want %#v", stub.record, want)
	}
}

func TestAuthServiceRegisterUsesBcryptCostAndUUID(t *testing.T) {
	stub := &authDataStub{}
	authService := newAuthServiceForTest(t, stub)

	err := authService.Register(context.Background(), model.PasswordRegisterReq{
		Username: "testuser",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	cost, err := bcrypt.Cost([]byte(stub.record.Password))
	if err != nil {
		t.Fatalf("stored password is not a bcrypt hash: %v", err)
	}
	if cost != passwordHashCost {
		t.Fatalf("stored bcrypt cost = %d, want %d", cost, passwordHashCost)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stub.record.Password), []byte("123456")); err != nil {
		t.Fatalf("stored bcrypt hash does not match password: %v", err)
	}
	if _, err := uuid.Parse(stub.record.UserID); err != nil {
		t.Fatalf("generated user_id = %q, want UUID: %v", stub.record.UserID, err)
	}
}

func TestAuthServiceRegisterRejectsInvalidInput(t *testing.T) {
	tests := []model.PasswordRegisterReq{
		{Username: "   ", Password: "123456"},
		{Username: strings.Repeat("用", 51), Password: "123456"},
		{Username: "testuser", Password: "12345"},
	}
	for _, req := range tests {
		stub := &authDataStub{}
		authService := newAuthServiceForTest(t, stub)

		err := authService.Register(context.Background(), req)
		if !errors.Is(err, ErrInvalidRegisterParams) {
			t.Fatalf("Register(%#v) error = %v, want ErrInvalidRegisterParams", req, err)
		}
		if stub.checkedUsername != "" || stub.createCalls != 0 {
			t.Fatalf("invalid request reached data layer: %#v", stub)
		}
	}
}

func TestAuthServiceRegisterReturnsUsernameExists(t *testing.T) {
	stub := &authDataStub{exists: true}
	authService := newAuthServiceForTest(t, stub)

	err := authService.Register(context.Background(), model.PasswordRegisterReq{
		Username: "testuser",
		Password: "123456",
	})
	if !errors.Is(err, ErrUsernameAlreadyExists) {
		t.Fatalf("Register() error = %v, want ErrUsernameAlreadyExists", err)
	}
	if stub.createCalls != 0 {
		t.Fatalf("CreateUser calls = %d, want 0", stub.createCalls)
	}
}

func TestAuthServiceRegisterMapsConcurrentDuplicate(t *testing.T) {
	stub := &authDataStub{createErr: data.ErrUsernameAlreadyExists}
	authService := newAuthServiceForTest(t, stub)
	authService.hashPassword = func([]byte, int) ([]byte, error) {
		return []byte("bcrypt-hash"), nil
	}
	authService.newUserID = func() string { return "user-uuid" }

	err := authService.Register(context.Background(), model.PasswordRegisterReq{
		Username: "testuser",
		Password: "123456",
	})
	if !errors.Is(err, ErrUsernameAlreadyExists) {
		t.Fatalf("Register() error = %v, want ErrUsernameAlreadyExists", err)
	}
}

func TestAuthServiceLoginReturnsTokenAndCamelCaseUserData(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate test password hash: %v", err)
	}
	stub := &authDataStub{user: data.UserRecord{
		UserID:     "abc123",
		Username:   "testuser",
		Password:   string(passwordHash),
		NickName:   "测试用户",
		Avatar:     "https://example.com/avatar.jpg",
		Sex:        1,
		VIP:        false,
		Phone:      "",
		Email:      "",
		UserStatus: 1,
		UserType:   1,
	}}
	authService := newAuthServiceForTest(t, stub)
	now := time.Date(2030, time.January, 2, 3, 4, 5, 0, time.UTC)
	authService.now = func() time.Time { return now }

	got, err := authService.Login(context.Background(), model.PasswordLoginReq{
		Username: "  testuser  ",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if stub.loginUsername != "testuser" {
		t.Fatalf("queried username = %q, want testuser", stub.loginUsername)
	}
	wantUserInfo := model.PasswordLoginUserInfo{
		UserID: "abc123", Username: "testuser", NickName: "测试用户",
		Avatar: "https://example.com/avatar.jpg", Sex: 1, VIP: false,
		Phone: "", Email: "", UserStatus: 1, UserType: 1,
	}
	if got.UserInfo != wantUserInfo {
		t.Fatalf("userInfo = %#v, want %#v", got.UserInfo, wantUserInfo)
	}

	parsedToken, err := jwt.ParseWithClaims(
		got.Token,
		&authClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				t.Fatalf("signing method = %v, want HS256", token.Method.Alg())
			}
			return []byte("test-jwt-secret"), nil
		},
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil || !parsedToken.Valid {
		t.Fatalf("parse generated token: valid=%t error=%v", parsedToken != nil && parsedToken.Valid, err)
	}
	claims, ok := parsedToken.Claims.(*authClaims)
	if !ok {
		t.Fatalf("token claims type = %T", parsedToken.Claims)
	}
	if claims.UserID != "abc123" {
		t.Fatalf("token user_id = %q, want abc123", claims.UserID)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.Time.Equal(now.Add(24*time.Hour)) {
		t.Fatalf("token expires_at = %v, want %v", claims.ExpiresAt, now.Add(24*time.Hour))
	}
}

func TestAuthServiceLoginRejectsUnknownUserWrongPasswordAndDisabledUser(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate test password hash: %v", err)
	}
	tests := []struct {
		name     string
		stub     *authDataStub
		password string
	}{
		{name: "unknown user", stub: &authDataStub{userErr: data.ErrUserNotFound}, password: "123456"},
		{name: "wrong password", stub: &authDataStub{user: data.UserRecord{Password: string(passwordHash), UserStatus: 1}}, password: "654321"},
		{name: "disabled user", stub: &authDataStub{user: data.UserRecord{Password: string(passwordHash), UserStatus: -1}}, password: "123456"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authService := newAuthServiceForTest(t, test.stub)
			_, err := authService.Login(context.Background(), model.PasswordLoginReq{
				Username: "testuser",
				Password: test.password,
			})
			if !errors.Is(err, ErrInvalidCredentials) {
				t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
			}
		})
	}
}

func TestAuthServiceLoginRejectsInvalidInput(t *testing.T) {
	stub := &authDataStub{}
	authService := newAuthServiceForTest(t, stub)
	_, err := authService.Login(context.Background(), model.PasswordLoginReq{
		Username: " ",
		Password: "12345",
	})
	if !errors.Is(err, ErrInvalidLoginParams) {
		t.Fatalf("Login() error = %v, want ErrInvalidLoginParams", err)
	}
	if stub.loginUsername != "" {
		t.Fatalf("invalid login reached data layer with username %q", stub.loginUsername)
	}
}

func TestNewAuthServiceRejectsInvalidJWTConfig(t *testing.T) {
	stub := &authDataStub{}
	if _, err := NewAuthService(stub, config.JWTConfig{Expire: 24}); err == nil {
		t.Fatal("NewAuthService() with empty secret error = nil")
	}
	if _, err := NewAuthService(stub, config.JWTConfig{Secret: "secret", Expire: 0}); err == nil {
		t.Fatal("NewAuthService() with zero expiration error = nil")
	}
	if _, err := NewAuthService(stub, config.JWTConfig{
		Secret:           "secret",
		Expire:           24,
		TokenCacheExpire: 12,
	}); err == nil {
		t.Fatal("NewAuthService() with token cache shorter than token expiration error = nil")
	}
}

func TestAuthServiceLogoutUsesConfiguredTokenCacheTTL(t *testing.T) {
	stub := &authDataStub{}
	authService := newAuthServiceForTest(t, stub)
	now := time.Date(2030, time.January, 2, 3, 4, 5, 0, time.UTC)
	authService.now = func() time.Time { return now }
	token, err := authService.signToken("abc123")
	if err != nil {
		t.Fatalf("signToken() error = %v", err)
	}

	if err := authService.Logout(context.Background(), token); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if stub.blacklistToken != token {
		t.Fatalf("blacklisted token = %q, want generated token", stub.blacklistToken)
	}
	if stub.blacklistTTL != 48*time.Hour {
		t.Fatalf("blacklist TTL = %s, want 48h", stub.blacklistTTL)
	}
}

func TestAuthServiceLogoutRejectsInvalidOrExpiredToken(t *testing.T) {
	stub := &authDataStub{}
	authService := newAuthServiceForTest(t, stub)
	now := time.Date(2030, time.January, 2, 3, 4, 5, 0, time.UTC)
	authService.now = func() time.Time { return now }

	for _, token := range []string{"not-a-jwt", ""} {
		if err := authService.Logout(context.Background(), token); !errors.Is(err, ErrInvalidLogoutToken) {
			t.Fatalf("Logout(%q) error = %v, want ErrInvalidLogoutToken", token, err)
		}
	}

	authService.jwtExpire = -time.Hour
	expiredToken, err := authService.signToken("abc123")
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	if err := authService.Logout(context.Background(), expiredToken); !errors.Is(err, ErrInvalidLogoutToken) {
		t.Fatalf("Logout(expired token) error = %v, want ErrInvalidLogoutToken", err)
	}
	if stub.blacklistToken != "" {
		t.Fatalf("expired token was blacklisted: %q", stub.blacklistToken)
	}
}

func TestAuthServiceLogoutReturnsBlacklistError(t *testing.T) {
	stub := &authDataStub{blacklistErr: errors.New("redis unavailable")}
	authService := newAuthServiceForTest(t, stub)
	now := time.Date(2030, time.January, 2, 3, 4, 5, 0, time.UTC)
	authService.now = func() time.Time { return now }
	token, err := authService.signToken("abc123")
	if err != nil {
		t.Fatalf("signToken() error = %v", err)
	}

	if err := authService.Logout(context.Background(), token); err == nil {
		t.Fatal("Logout() error = nil, want blacklist error")
	}
}
