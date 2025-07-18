package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
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
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	PubDate     rfc1223zTime `xml:"pubDate"`
}

type rfc1223zTime time.Time

func (c *rfc1223zTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC1123Z, v)
	if err != nil {
		return err
	}
	*c = rfc1223zTime(parse)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ret := RSSFeed{}
	if err := xml.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	ret.Channel.Title = html.UnescapeString(ret.Channel.Title)
	ret.Channel.Description = html.UnescapeString(ret.Channel.Description)

	for i, item := range ret.Channel.Item {
		ret.Channel.Item[i].Title = html.UnescapeString(item.Title)
		ret.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &ret, nil
}
