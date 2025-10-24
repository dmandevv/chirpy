package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func makeHMACToken(subject string, secret string, exp time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(exp.UTC()),
		Subject:   subject,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

func TestValidateJWT_ValidToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	id := uuid.New()
	token, err := makeHMACToken(id.String(), secret, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error for valid token: %v", err)
	}
	if got != id {
		t.Fatalf("expected %v, got %v", id, got)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	badSecret := "other-secret"
	id := uuid.New()
	token, err := makeHMACToken(id.String(), secret, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	got, err := ValidateJWT(token, badSecret)
	if err == nil {
		t.Fatalf("expected error when validating with wrong secret, got nil and uuid %v", got)
	}
	if got != uuid.Nil {
		t.Fatalf("expected uuid.Nil when validation fails, got %v", got)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	id := uuid.New()
	token, err := makeHMACToken(id.String(), secret, time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil and uuid %v", got)
	}
	// jwt/v5 returns jwt.ErrTokenExpired for expired tokens; ensure it's reported (or at least an error)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		// still accept any non-nil error, but note if it's not the specific expiration error
		t.Logf("expected jwt.ErrTokenExpired (or related) but got: %v", err)
	}
	if got != uuid.Nil {
		t.Fatalf("expected uuid.Nil for expired token, got %v", got)
	}
}

func TestValidateJWT_InvalidSubject(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	token, err := makeHMACToken("not-a-uuid", secret, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("expected error for token with invalid subject, got nil and uuid %v", got)
	}
	if got != uuid.Nil {
		t.Fatalf("expected uuid.Nil when subject is invalid, got %v", got)
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	got, err := ValidateJWT("not.a.valid.token", secret)
	if err == nil {
		t.Fatalf("expected error for malformed token, got nil and uuid %v", got)
	}
	if got != uuid.Nil {
		t.Fatalf("expected uuid.Nil for malformed token, got %v", got)
	}
}
func TestGetBearerToken_Valid(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Add("Authorization", "Bearer abc123.def456.ghi789")

	got, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error for valid header: %v", err)
	}
	if want := "abc123.def456.ghi789"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	t.Parallel()

	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for missing Authorization header, got nil")
	}
}

func TestGetBearerToken_NotBearer(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Add("Authorization", "Basic abc123")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for non-Bearer Authorization header, got nil")
	}
}

func TestGetBearerToken_EmptyToken(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Add("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
