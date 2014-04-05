package command

import (
	"github.com/codegangsta/cli"
)

func StatusCommand() cli.Command {
	return cli.Command{
		Name:  "status",
		Usage: "get status",
		Action: func(c *cli.Context) {
			env := getCommandEnv(c)
			defer env.DB.Close()
			StatusAction(env, c)
		},
	}
}

func IndexCommand() cli.Command {
	return cli.Command{
		Name:  "index",
		Usage: "add path to index",
		Action: func(c *cli.Context) {
			env := getCommandEnv(c)
			defer env.DB.Close()
			IndexAction(env, c)
		},
	}
}

func getCommandEnv(c *cli.Context) *Environment {
	return GetEnvironment(c.GlobalString("dir"), c.GlobalInt("workers"), c.GlobalBool("verbose"))
}
