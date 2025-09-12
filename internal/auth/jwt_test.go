package auth_test

import (
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
