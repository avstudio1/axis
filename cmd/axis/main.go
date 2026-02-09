package main

import (
	"context"
	"log"
	"os"

	"axis/internal/workspace"

	"github.com/joho/godotenv"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

func main() {
	// 1. Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on shell variables")
	}

	ctx := context.Background()

	adminEmail := os.Getenv("ADMIN_EMAIL")
	serviceAccountEmail := os.Getenv("SERVICE_ACCOUNT_EMAIL")

	if adminEmail == "" || serviceAccountEmail == "" {
		log.Fatal("Required environment variables are missing")
	}

	ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: serviceAccountEmail,
		Subject:         adminEmail,
		Scopes: []string{
			"https://www.googleapis.com/auth/admin.directory.user",
		},
	})
	if err != nil {
		log.Fatalf("Token source error: %v", err)
	}

	adminSvc, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		log.Fatalf("Service init error: %v", err)
	}

	ws := workspace.NewService(adminSvc)

	testEmail := os.Getenv("TEST_USER_EMAIL")
	user, err := ws.GetUser(testEmail)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	log.Printf("Verified: %s (%s)", user.Name, user.Email)
}
