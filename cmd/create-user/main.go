package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/joho/godotenv"
)

func main() {
	email := flag.String("email", "", "User email")
	password := flag.String("password", "", "User password")
	username := flag.String("username", "", "Username")
	flag.Parse()

	if *email == "" || *password == "" {
		log.Fatal("email and password are required")
	}

	if err := godotenv.Load(); err != nil {
		log.Printf(".env was not loaded, using system environment variables: %v", err)
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(clerkSecretKey)

	ctx := context.Background()

	existingUsers, err := clerkuser.List(ctx, &clerkuser.ListParams{
		EmailAddresses: []string{*email},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(existingUsers.Users) > 0 {
		fmt.Printf("User already exists in Clerk: %s\n", existingUsers.Users[0].ID)
		return
	}

	emailAddresses := []string{*email}
	params := &clerkuser.CreateParams{
		EmailAddresses: &emailAddresses,
		Password:       password,
	}
	if *username != "" {
		params.Username = username
	}

	createdUser, err := clerkuser.Create(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Clerk user: %s\n", createdUser.ID)
}
