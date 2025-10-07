package main

//import _ "github.com/lib/pq"
import (
    "os"
    "log"
    "database/sql"

    "github.com/kevin-baik/aggreGator/internal/config"
    "github.com/kevin-baik/aggreGator/internal/database"
    _ "github.com/lib/pq"
)

type state struct {
    db	    *database.Queries
    config  *config.Config
}

func main() {
    cfg, err := config.Read()
    if err != nil {
	log.Fatalf("Error reading config file: %v", err)
    }
    dbURL := cfg.DBUrl
    db, err := sql.Open("postgres", dbURL) 
    if err != nil {
	log.Fatalf("Error opening db @ %v", dbURL)
    }
    dbQueries := database.New(db)

    programState := &state{
	db:	dbQueries,
	config: &cfg,
    }
    
    cmds := commands{
	registeredCommands: make(map[string]func(*state, command) error),
    }
    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("users", handlerListUsers)
    cmds.register("reset", handlerResetDatabase)
    cmds.register("agg", handlerAgg)
    cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    cmds.register("feeds", handlerAllFeeds)
    cmds.register("follow", middlewareLoggedIn(handlerFollow))
    cmds.register("following", middlewareLoggedIn(handlerFollowing))
    cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

    if len(os.Args) < 2 {
	log.Fatalf("No command provided... Usage: cli <command> [args...]")	
    }

    cmdName := os.Args[1]
    cmdArgs := os.Args[2:]

    cmd := command{
	Name: cmdName,
	Args: cmdArgs,
    }
    if err := cmds.run(programState, cmd); err != nil {
	log.Fatalf("Command Run Error:", err)
    }
}
