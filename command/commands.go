package command

import (
  "github.com/codegangsta/cli"
)

func StatusCommand() cli.Command {
  return cli.Command{
    Name:   "status",
    Usage:  "get status",
    Action: func(c *cli.Context) {
      env := GetEnvironment(c.String("dir"))
      StatusAction(env, c)
    },
  }
}

func IndexCommand() cli.Command {
  return cli.Command{
    Name: "index",
    Usage: "add path to index",
    Action: func(c *cli.Context) {
      env := GetEnvironment(c.String("dir"))
      IndexAction(env, c)
    },
  }
}
