package auth

import (
	"testing"
	"time"
)

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("s3cret-password")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "s3cret-password" {
		t.Fatal("hash must not equal plaintext")
	}
	if !VerifyPassword(hash, "s3cret-password") {
		t.Fatal("expected password to verify")
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Fatal("expected wrong password to fail")
	}
}

func TestTokenIssueAndParse(t *testing.T) {
	svc := NewTokenService("test-secret-0123456789abcdef", time.Hour)

	token, err := svc.Issue(42, "admin@test.dev")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
	if claims.Email != "admin@test.dev" {
		t.Fatalf("expected email admin@test.dev, got %q", claims.Email)
	}
}

func TestTokenRejectsWrongSecret(t *testing.T) {
	issuer := NewTokenService("secret-one", time.Hour)
	verifier := NewTokenService("secret-two", time.Hour)

	token, err := issuer.Issue(1, "a@b.dev")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verifier.Parse(token); err == nil {
		t.Fatal("expected token signed with another secret to be rejected")
	}
}

func TestTokenRejectsTampering(t *testing.T) {
	svc := NewTokenService("test-secret-0123456789abcdef", time.Hour)

	token, err := svc.Issue(1, "a@b.dev")
	if err != nil {
		t.Fatal(err)
	}
	tampered := token[:len(token)-4] + "AAAA"
	if _, err := svc.Parse(tampered); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestTokenExpired(t *testing.T) {
	svc := NewTokenService("test-secret-0123456789abcdef", -time.Hour)

	token, err := svc.Issue(1, "a@b.dev")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Parse(token); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}
