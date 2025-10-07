package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    "log"
    "strings"

    "github.com/kevin-baik/aggreGator/internal/database"
    "github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("Time duration not provided. Usage: agg <time_duration>)\n")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: agg <time_duration>)\n")
    }
    
    timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
    if err != nil {
	return fmt.Errorf("Unable to parse time duration")
    }
    log.Printf("Collecting feeds every: %v", timeBetweenRequests)

    ticker := time.NewTicker(timeBetweenRequests)

    for ; ; <-ticker.C {
	scrapeFeeds(s)
    }

    return nil
}

func scrapeFeeds(s *state) {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
	log.Println("Couldn't get next feed to fetch: %w", err)
	return
    }
    log.Println("Found a feed to fetch!")
    scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
    _, err := db.MarkFeedFetched(context.Background(), feed.ID)
    if err != nil {
	log.Printf("Couldn't mark %v: %w", feed.Name, err)
	return
    }
    
    rssFeed, err := fetchFeed(context.Background(), feed.Url)
    if err != nil {
	log.Printf("Couldn't collect feed %v: %w", feed.Name, err)
	return
    }
    
    for _, item := range rssFeed.Channel.Item {
	publishedAt := sql.NullTime{}
	if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
	    publishedAt = sql.NullTime{
		Time:  t,
		Valid: true,
	    }
	}

	_, err := db.CreatePost(context.Background(), database.CreatePostParams{
	    ID:          uuid.New(),
	    CreatedAt:   time.Now().UTC(),
	    UpdatedAt:   time.Now().UTC(),
	    Title:       item.Title,
	    Url:         item.Link,
	    Description: sql. NullString{
		String:	item.Description,
		Valid:	true,
	    },
	    PublishedAt: publishedAt,
	    FeedID:      feed.ID,
	})
	if err != nil {
	    if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		continue
	    }
	    log.Printf("Couldn't create post: %v", err)
	    continue
	}
    }
    log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
