package main

import(
    "net/http"
    "context"
    "fmt"
    "io"
    "time"
    "encoding/xml"
    "html"

    "github.com/kevin-baik/aggreGator/internal/database"
    "github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    client := &http.Client{}
    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
	fmt.Println("Error creating request")
	return &RSSFeed{}, err
    }
    
    req.Header.Set("User-Agent", "aggreGator")

    resp, err := client.Do(req)
    if err != nil {
	fmt.Println("Error making request")
	return &RSSFeed{}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
	fmt.Println("Error reading response body")
	return &RSSFeed{}, err
    }

    var rssFeed RSSFeed
    err = xml.Unmarshal(body, &rssFeed)
    if err != nil {
	return &RSSFeed{}, err
    }
    
    rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
    rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
    
    for _, item := range rssFeed.Channel.Item {
	item.Title = html.UnescapeString(item.Title)
	item.Description = html.UnescapeString(item.Description)
    }

    return &rssFeed, nil
}

func handlerFetchFeed(s *state, cmd command) error {
    feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
	return err
    }

    fmt.Println(feed)
    return nil
}

func handlerAddFeed(s *state, cmd command) error {
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

    user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
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

func handlerFollow(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
	return fmt.Errorf("URL to follow not provided. Usage: follow (feed_url)\n")
    } else if len(cmd.Args) > 1 {
	return fmt.Errorf("Too many arguments. Usage: follow (feed_url)\n")
    }

    feedUrl := cmd.Args[0]
    user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
    if err != nil {
	return err
    }
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
	return err
    }

    fmt.Println("Feed:", feed.Name)
    fmt.Println("User:", user.Name)

    return nil
}

func handlerFollowing(s *state, cmd command) error {
    if len(cmd.Args) > 0 {
	return fmt.Errorf("No arguments needed. Usage: following\n")
    }

    user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
    if err != nil {
	return err
    }
    
    feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
	return err
    }

    for i, feedFollow := range feedFollows {
	fmt.Printf("---------- FEED #%v ----------\n", i + 1)
	fmt.Println("Feed:", feedFollow.FeedName)
	fmt.Println("Name:", feedFollow.UserName)
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
