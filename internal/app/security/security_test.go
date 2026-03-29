package security

import (
	"testing"
)

func TestOpaqueTokenHashStable(t *testing.T) {
	token := NewOpaqueToken()
	if token == "" {
		t.Fatal("token should not be empty")
	}
	hash := HashOpaqueToken(token)
	if hash != HashOpaqueToken(token) {
		t.Fatal("hash should be stable")
	}
}

func TestAccessTokenRoundTrip(t *testing.T) {
	token, _, err := SignAccessToken(1001, "admin")
	if err != nil {
		t.Fatalf("sign access token: %v", err)
	}
	claims, err := ParseAccessToken(token)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if claims.UID != 1001 || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
