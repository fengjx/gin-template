package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"gin-template/internal/app/config"
	"gin-template/internal/app/trace"
	"gin-template/pkg/errs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccessClaims struct {
	UID  int64  `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.Wrap(err, "生成密码哈希失败")
	}
	return string(hash), nil
}

func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errs.Wrap(err, "校验密码失败")
	}
	return nil
}

func SignAccessToken(uid int64, role string) (string, time.Time, error) {
	cfg := config.Get()
	expireAt := time.Now().Add(time.Duration(cfg.Auth.AccessTokenTTLMinutes) * time.Minute)
	subject := strconv.FormatInt(uid, 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Auth.Issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	})
	signed, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		return "", time.Time{}, errs.Wrap(err, "签发访问令牌失败")
	}
	return signed, expireAt, nil
}

func ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(config.Get().Auth.JWTSecret), nil
	})
	if err != nil {
		return nil, errs.Wrap(err, "解析访问令牌失败")
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errs.WithStack(errors.New("invalid access token"))
	}
	return claims, nil
}

func NewOpaqueToken() string {
	return uuid.NewString() + "." + uuid.NewString()
}

func HashOpaqueToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func MaskToken(token string) string {
	if len(token) < 8 {
		return "****"
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

func TokenContext(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		traceID = trace.Generate()
	}
	return trace.WithTraceID(ctx, traceID)
}
