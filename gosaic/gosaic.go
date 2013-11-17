package main

import (
        "flag"
        "github.com/atongen/gosaic/runner"
        "log"
        "os"
)

var path string

func init() {
        var defaultPath string
        var err error
        defaultPath, err = os.Getwd()
        if err != nil {
                log.Fatalf("Unable to get current directory.")
        }
        flag.StringVar(&path, "p", defaultPath, "Project path")
}

func main() {
        var run runner.Runner

        flag.Parse()
        subcommand := flag.Arg(0)
        arg := flag.Arg(1)

        switch subcommand {
        case "init":
          run = runner.Init{Path: path, Arg: arg}
        default:
                log.Fatalf("Invalid sub-command: %s\n", subcommand)
        }

        err := run.Execute()
        if err != nil {
                log.Fatal(err)
        }
}
