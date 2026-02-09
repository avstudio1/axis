package workspace

import (
	"testing"
)

func TestMockService(t *testing.T) {
	mock := &MockService{
		OnGetUser: func(email string) (*User, error) {
			return &User{
				Email: "test@axis.com",
				Name:  "Test User",
			}, nil
		},
	}

	user, err := mock.GetUser("test@axis.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "Test User" {
		t.Errorf("expected 'Test User', got %s", user.Name)
	}
}
