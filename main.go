package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tchoffman/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

// Create a login handler function: func handlerLogin(s *state, cmd command) error. This will be the function signature of all command handlers.

// If the command's arg's slice is empty, return an error; the login handler expects a single argument, the username.
// Use the state's access to the config struct to set the user to the given username. Remember to return any errors.
// Print a message to the terminal that the user has been set.

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects a single argument, the username")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}

	fmt.Printf("User set to %s\n", cmd.args[0])
	return nil
}

// Create a commands struct. This will hold all the commands the CLI can handle.
// Add a map[string]func(*state, command) error field to it.
// This will be a map of command names to their handler functions.

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return f(s, cmd)
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{cfg: &cfg}

	// Create a new instance of the commands struct with an initialized map of handler functions.
	cmds := commands{cmds: make(map[string]func(*state, command) error)}

	// Register a handler function for the login command.
	cmds.register("login", handlerLogin)

	// Use os.Args to get the command-line arguments passed in by the user.
	// If there are no arguments, print a message to the terminal and exit.
	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		os.Exit(1)
	}

	// Parse the command-line arguments into a command struct.
	// The first argument is the command name, and the rest are the arguments to the command.
	cmd := command{name: os.Args[1], args: os.Args[2:]}
	fmt.Printf("Running command: %+v\n", cmd)
	// Use the commands.run method to run the given command and print any errors returned.
	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}

}
