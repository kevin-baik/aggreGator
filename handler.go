package main

import (
    "fmt"
    "context"
    "time"
    
    "github.com/kevin-baik/aggreGator/internal/database"
    "github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("No login user provided")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: login [arg]\n")
    }
    name := cmd.Args[0]
    user, err := s.db.GetUser(context.Background(), name)
    if err != nil {
	return fmt.Errorf("Login failed. User not in system")
    }

    if err := s.config.SetUser(user.Name); err != nil {
	return err 
    }
    fmt.Println(user.Name, "has been set as current user")
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("No registration user provided")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: register [arg]\n")
    }
    name := cmd.Args[0]
    user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
	ID:	    uuid.New(),
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
	Name:	    name,
    })
    if err != nil {
	return fmt.Errorf("Create User Error: %w", err)
    }

    if err := s.config.SetUser(user.Name); err != nil {
	return err
    }
    fmt.Printf("Created User: %v\n", user.Name)
    printUser(user)
    return nil
}

func handlerListUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(context.Background())
    if err != nil {
	return err
    }

    for _, user := range users {
	if user.Name == s.config.CurrentUserName {
	    fmt.Println("*", user.Name, "(current)")
	} else {
	    fmt.Println("*", user.Name)
	}
    }
    return nil
}

func handlerResetDatabase(s *state, cmd command) error {
    if err := s.db.DeleteUsers(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    if err := s.db.DeleteFeeds(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    if err := s.db.DeleteFeedFollows(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    fmt.Println("Database reset successful!")
    return nil
}

func printUser(user database.User) {
    fmt.Printf(" * ID:      %v\n", user.ID)
    fmt.Printf(" * Name:    %v\n", user.Name)
}
