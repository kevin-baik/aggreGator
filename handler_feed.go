package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kevin-baik/aggreGator/internal/database"
    "github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("Feed name not provided. Usage: addfeed (feed_name, feed_url)\n")
    } else if len(cmd.Args) == 1 {
	return fmt.Errorf("Feed URL not provided. Usage: addfeed (feed_name, feed_url)\n")
    } else if len(cmd.Args) > 2 {
	return fmt.Errorf("Too many arguments. Usage: addfeed (feed_name, feed_url)\n")
    }

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
	ID:	    uuid.New(),
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
	Name:	    cmd.Args[0],
	Url:	    cmd.Args[1],
    })
        if err != nil {
	fmt.Println("Error creating feed of: %v", feed.Name)
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
	return err
    }

    printFeed(feed)
    return nil
}

func handlerAllFeeds(s *state, cmd command) error {
    if len(cmd.Args) > 0 {
	return fmt.Errorf("Too many arguments. Usage: feed\n")
    }

    feeds, err := s.db.GetAllFeeds(context.Background())
    if err != nil {
	return fmt.Errorf("Error getting all feeds: %w", err)
    }

    for i, feed := range feeds {
	fmt.Printf("---------- FEED #%v ----------\n", i + 1)
	printFeed(feed)
	fmt.Println("-----------------------------")
    }
    return nil
}

func printFeed(feed database.Feed) {
    fmt.Printf(" * ID:        %v\n", feed.ID)
    fmt.Printf(" * CreatedAt: %v\n", feed.CreatedAt)
    fmt.Printf(" * UpdatedAt: %v\n", feed.UpdatedAt)
    fmt.Printf(" * Name:      %v\n", feed.Name)
    fmt.Printf(" * URL:       %v\n", feed.Url)
}
