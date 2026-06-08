package main

import (
	"context"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

type fakeOAuth2Config struct {
	token *oauth2.Token
	err   error
	calls int
}

func (f *fakeOAuth2Config) Token(_ context.Context) (*oauth2.Token, error) {
	f.calls++
	return f.token, f.err
}

func TestOauthbearerTokenProviderCacheHit(t *testing.T) {
	t.Parallel()

	fake := &fakeOAuth2Config{
		token: &oauth2.Token{
			AccessToken: "token-1",
			Expiry:      time.Now().Add(time.Hour),
		},
	}
	provider := newOauthbearerTokenProvider(fake)

	first, err := provider.Token()
	if err != nil {
		t.Fatalf("first Token() error: %v", err)
	}
	if first.Token != "token-1" {
		t.Fatalf("first token = %q, want %q", first.Token, "token-1")
	}
	if fake.calls != 1 {
		t.Fatalf("Token() calls = %d, want 1", fake.calls)
	}

	second, err := provider.Token()
	if err != nil {
		t.Fatalf("second Token() error: %v", err)
	}
	if second.Token != "token-1" {
		t.Fatalf("second token = %q, want %q", second.Token, "token-1")
	}
	if fake.calls != 1 {
		t.Fatalf("Token() calls after cache hit = %d, want 1", fake.calls)
	}
}

func TestOauthbearerTokenProviderCacheMissOnExpiry(t *testing.T) {
	t.Parallel()

	fake := &fakeOAuth2Config{
		token: &oauth2.Token{
			AccessToken: "token-1",
			Expiry:      time.Now().Add(-time.Minute),
		},
	}
	provider := newOauthbearerTokenProvider(fake)

	_, err := provider.Token()
	if err != nil {
		t.Fatalf("first Token() error: %v", err)
	}
	if fake.calls != 1 {
		t.Fatalf("Token() calls = %d, want 1", fake.calls)
	}

	fake.token = &oauth2.Token{
		AccessToken: "token-2",
		Expiry:      time.Now().Add(time.Hour),
	}

	got, err := provider.Token()
	if err != nil {
		t.Fatalf("second Token() error: %v", err)
	}
	if got.Token != "token-2" {
		t.Fatalf("token = %q, want %q", got.Token, "token-2")
	}
	if fake.calls != 2 {
		t.Fatalf("Token() calls after refresh = %d, want 2", fake.calls)
	}
}
