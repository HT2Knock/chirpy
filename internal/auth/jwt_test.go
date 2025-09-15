package auth_test

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/T2Knock/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	jwt, _ := auth.MakeJWT(userID, "thetoken", time.Hour)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        string
		wantErr     bool
	}{
		{
			name:        "Correct JWT",
			userID:      userID,
			tokenSecret: "thetoken",
			expiresIn:   time.Hour,
			want:        jwt,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("MakeJWT() = %v, want %v", got, tt.want)
			}

			uuid, err := auth.ValidateJWT(jwt, "thetoken")
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", err)
				}
				return
			}

			if uuid != tt.userID {
				t.Errorf("Validate token failed not match userID")
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJkOWYwZDU4Mi0yMzkzLTQ5YWQtOGY1NC1mMjY5Y2RkZTJjOWEiLCJleHAiOjE3NTc4Mzk1NzQsImlhdCI6MTc1NzgzNTk3NH0.083qQaVdYXF29vkTmgdpQ02VY8nphvtFPPfOZfJDF6w"
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "Valid bearer token extract",
			headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
			},
			want:    token,
			wantErr: false,
		},
		{
			name:    "Missing authorization header",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name: "Missing authorization token",
			headers: http.Header{
				"Authorization": []string{"Bearer"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Wrong scheme",
			headers: http.Header{
				"Authorization": []string{"Basic 123"},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token, err := auth.MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken() returned error: %v", err)
	}

	if len(token) != 64 {
		t.Errorf("unexpected token length: got %d, want 64", len(token))
	}

	_, decodeErr := hex.DecodeString(token)
	if decodeErr != nil {
		t.Errorf("token is not valid hex: %v", decodeErr)
	}
}
