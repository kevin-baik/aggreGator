package main

import (
    "context"
    "fmt"
)

func handlerResetDatabase(s *state, cmd command) error {
    if err := s.db.DeleteAllUsers(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    if err := s.db.DeleteAllFeeds(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    if err := s.db.DeleteAllFeedFollows(context.Background()); err != nil {
	return fmt.Errorf("Unable to delete users: %w", err)
    }
    fmt.Println("Database reset successful!")
    return nil
}
