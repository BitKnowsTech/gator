package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bitknowstech/gator/internal/config"
	"github.com/bitknowstech/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

func main() {
	userConfig := config.Read()
	db, err := sql.Open("postgres", userConfig.DbUrl)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	dbQueries := database.New(db)

	programState := state{
		conf: &userConfig,
		db:   dbQueries,
	}

	cmds := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	cmdName, cmdArgs := os.Args[1], os.Args[2:] // if this is an issue it can be changed. Magic numbers: 0: program name, 1: command name, "2:": everything after command name

	err = cmds.run(&programState, command{name: cmdName, args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
