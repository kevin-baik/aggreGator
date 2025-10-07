package main

import(
    "net/http"
    "context"
    "io"
    "time"
    "encoding/xml"
    "html"
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
    client := &http.Client{
	Timeout: 10 * time.Second,
    }
    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
	return nil, err
    }
    
    req.Header.Set("User-Agent", "aggreGator")

    resp, err := client.Do(req)
    if err != nil {
	return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
	return nil, err
    }

    var rssFeed RSSFeed
    err = xml.Unmarshal(body, &rssFeed)
    if err != nil {
	return nil, err
    }
    
    rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
    rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
    
    for _, item := range rssFeed.Channel.Item {
	item.Title = html.UnescapeString(item.Title)
	item.Description = html.UnescapeString(item.Description)
    }

    return &rssFeed, nil
}
