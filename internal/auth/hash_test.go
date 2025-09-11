package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "thepassword"
	hashPassword1, _ := HashPassword(password1)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hashPassword1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "thepasswor",
			hash:     hashPassword1,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := CheckPasswordHash(tt.password, tt.hash)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckPasswordHash() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("CheckPasswordHash() succeeded unexpectedly")
			}
		})
	}
}
