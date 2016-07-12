package command

import (
	"github.com/atongen/gosaic/controller"
	"github.com/codegangsta/cli"
)

func Status() cli.Command {
	return cli.Command{
		Name:  "status",
		Usage: "get status",
		Action: func(c *cli.Context) {
			env := getCommandEnv(c)
			if !hasExpectedArgs(c.Args(), 0) {
				env.Fatalln("Unexpected arguments present.")
			}
			env.Init()
			defer env.Close()

			controller.Status(env)
		},
	}
}

func Index() cli.Command {
	return cli.Command{
		Name:  "index",
		Usage: "add path to index",
		Action: func(c *cli.Context) {
			env := getCommandEnv(c)
			if !hasExpectedArgs(c.Args(), 1) {
				env.Fatalln("Path argument is required.")
			}
			env.Init()
			defer env.Close()

			controller.Index(env, c.Args()[0])
		},
	}
}

func getCommandEnv(c *cli.Context) controller.Environment {
	return controller.GetProdEnv(c.GlobalString("dir"), c.GlobalInt("workers"), c.GlobalBool("verbose"), c.GlobalBool("debug"))
}

// hasExpectedArgs checks whether the number of args are as expected.
func hasExpectedArgs(args []string, expected int) bool {
	switch expected {
	case -1:
		if len(args) > 0 {
			return true
		} else {
			return false
		}
	default:
		if len(args) == expected {
			return true
		} else {
			return false
		}
	}
}
