package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	t.Parallel()

	const pw = "s3cr3tP@ssw0rd"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == pw {
		t.Fatal("hash must not equal plaintext")
	}

	ok, err := CheckPasswordHash(pw, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected password to match hash")
	}
}

func TestCheckPasswordHash_FailsForWrongPassword(t *testing.T) {
	t.Parallel()

	const pw = "s3cr3tP@ssw0rd"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	ok, err := CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if ok {
		t.Fatal("expected CheckPasswordHash to return false for incorrect password")
	}
}

func TestHashProducesDifferentValues(t *testing.T) {
	t.Parallel()

	const pw = "same-password"
	h1, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	h2, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if h1 == h2 {
		t.Fatal("expected different hashes for the same password (random salt)")
	}
}
