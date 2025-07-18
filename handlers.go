package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/bitknowstech/gator/internal/database"
	"github.com/google/uuid"
)

var errNotEnoughArgs = errors.New("not enough arguments for command")

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errNotEnoughArgs
	}

	username := cmd.args[0]

	usr, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	if err := s.conf.SetUser(usr.Name); err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", usr.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errNotEnoughArgs
	}

	username := cmd.args[0]

	passable := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	usr, err := s.db.CreateUser(context.Background(), passable)
	if err != nil {
		return err
	}
	s.conf.SetUser(usr.Name)
	fmt.Println("User", usr.Name, "created and logged in!")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	usrs, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		current := usr.Name == s.conf.CurrentUserName
		if current {
			fmt.Printf("* %s (current)\n", usr.Name)
			continue
		}
		fmt.Printf("* %s\n", usr.Name)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errNotEnoughArgs
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("* Added feed %s at URL %s\n", feed.Name, feed.Url)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feedrow := range feeds {
		fmt.Printf("--- %s ---\n* Url: %s\n* Added by: %s\n", feedrow.Name, feedrow.Url, feedrow.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errNotEnoughArgs
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("no such feed")
	}

	retRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("You followed: %s\n", retRow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("*", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errNotEnoughArgs
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errNotEnoughArgs
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return err
		}

		err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		cachedFeedUrls, err := s.db.GetPostUrlsByFeed(context.Background(), feed.ID)
		if err != nil {
			return err
		}

		rss, err := fetchFeed(context.Background(), feed.Url)
		if err != nil {
			return err
		}

		var postsToCreate []database.CreatePostParams

		for _, item := range rss.Channel.Item {
			if len(item.Title) < 1 {
				continue
			}
			if !slices.Contains(cachedFeedUrls, item.Link) {
				postsToCreate = append(postsToCreate, database.CreatePostParams{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Title:     item.Title,
					Url:       item.Link,
					Description: sql.NullString{
						String: item.Description,
						Valid:  len(item.Description) > 1,
					},
					PublishedAt: sql.NullTime{
						Time:  time.Time(item.PubDate),
						Valid: time.Time(item.PubDate) != time.Time{},
					},
					FeedID: feed.ID,
				})
			}
		}

		// Here we insert the post into the database. for post := range item...
		// this lets us run against the post for later
		fmt.Println("* Fetched", len(postsToCreate), "posts from feed", feed.Name)
		for _, item := range postsToCreate {
			_, err := s.db.CreatePost(context.Background(), item)
			if err != nil {
				return err
			}
		}
	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	posts, err := s.db.GetPostsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("*", "["+post.PublishedAt.Time.Format(time.DateOnly)+"]", post.Title, "["+post.Url+"]")
	}
	return nil
}

// resets the DB, should account for any command that contains "reset" in the name. e.g. ResetUsers/ResetPosts/ResetSubscriptions
func handlerReset(s *state, cmd command) error {
	err := s.db.ResetFeeds(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetFeedFollows(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetPosts(context.Background())
	if err != nil {
		return err
	}
	return nil
}
