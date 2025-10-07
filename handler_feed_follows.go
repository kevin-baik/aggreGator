package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kevin-baik/aggreGator/internal/database"
    "github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("URL to follow not provided. Usage: follow (feed_url)\n")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: follow (feed_url)\n")
    }

    feedUrl := cmd.Args[0]
    feed, err := s.db.GetFeedWithURL(context.Background(), feedUrl)
    if err != nil {
	return err
    }

    _, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
	ID:	    uuid.New(),
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
	UserID:	    user.ID,
	FeedID:	    feed.ID,
    })
    if err != nil {
	return fmt.Errorf("Unable to create feed follow: %w", err)
    }
    
    printFeedFollow(user.Name, feed.Name)
    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    if len(cmd.Args) > 0 {
	return fmt.Errorf("No arguments needed. Usage: following\n")
    }
    
    feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
	return err
    }

    fmt.Printf("---------- %s FEEDS ----------\n", user.Name)
    for i, feedFollow := range feedFollows {
	fmt.Printf("Feed #%v: %v\n", i + 1, feedFollow.FeedName)
    }
    fmt.Println("-----------------------------------")
    return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("URL to follow not provided. Usage: follow (feed_url)\n")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: follow (feed_url)\n")
    }

    feedUrl := cmd.Args[0]
    feed, err := s.db.GetFeedWithURL(context.Background(), feedUrl)
    if err != nil {
	return err
    }
    
    err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
	UserID:	    user.ID,
	FeedID:	    feed.ID,
    })
    if err != nil {
	return err
    }
    return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
