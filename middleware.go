package main

import (
	"context"
	"errors"

	"github.com/bitknowstech/gator/internal/database"
)

var errCouldntGetLoggedInUser = errors.New("couldn't get logged in user")

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			return errCouldntGetLoggedInUser
		}

		return handler(s, cmd, user)
	}
}
