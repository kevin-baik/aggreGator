package main

import (
    "context"
    "fmt"
    "time"
    "log"
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

func scrapeFeeds(s *state) error {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
	return fmt.Errorf("couldn't fetch next feed: %w", err)
    }
    
    _, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
    if err != nil {
	return fmt.Errorf("couldn't mark feed: %w", err)
    }
    
    rssFeed, err := fetchFeed(context.Background(), feed.Url)
    if err != nil {
	return fmt.Errorf("couldn't fetch feed: %w", err)
    }
    
    printRSSFeed(*rssFeed)
    return nil
}
