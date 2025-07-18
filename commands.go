package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("state not found")
	}

	cmdFunc, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	if err := cmdFunc(s, cmd); err != nil {
		return fmt.Errorf("error in command %s: %v", cmd.name, err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
