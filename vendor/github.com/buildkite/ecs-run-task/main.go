package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/buildkite/ecs-run-task/runner"
	"github.com/urfave/cli"
)

var (
	Version string
)

func main() {
	app := cli.NewApp()
	app.Name = "ecs-run-task"
	app.Usage = "run a once-off task on ECS and tail the output from cloudwatch"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Show debugging information",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Task definition file in JSON or YAML",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Task name",
		},
		cli.StringFlag{
			Name:  "cluster, c",
			Value: "default",
			Usage: "ECS cluster name",
		},
		cli.StringFlag{
			Name:  "log-group, l",
			Value: "ecs-task-runner",
			Usage: "Cloudwatch Log Group Name to write logs to",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		requireFlagValue(ctx, "file")

		if _, err := os.Stat(ctx.String("file")); err != nil {
			return cli.NewExitError(err, 1)
		}

		if !ctx.Bool("debug") {
			log.SetOutput(ioutil.Discard)
		}

		r := runner.New()
		r.TaskDefinitionFile = ctx.String("file")
		r.Cluster = ctx.String("cluster")
		r.TaskName = ctx.String("name")
		r.LogGroupName = ctx.String("log-group")

		if err := r.Run(); err != nil {
			if ec, ok := err.(cli.ExitCoder); ok {
				return ec
			}
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	}

	app.Run(os.Args)
}

func requireFlagValue(ctx *cli.Context, name string) {
	if ctx.String(name) == "" {
		fmt.Printf("ERROR: Required flag %q isn't set\n\n", name)
		cli.ShowAppHelpAndExit(ctx, 1)
	}
}
