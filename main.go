package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/buildkite/ecs-run-task/parser"
	"github.com/urfave/cli/v2"
)

const ECS_POLL_INTERVAL = 1 * time.Second
const version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "ecs-upload-task"
	app.Usage = "upload ecs task definitions and update ECS services"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "file", Value: "taskdefinition.json", Usage: "the task definition to upload"},
		&cli.StringFlag{Name: "cluster", Value: "default", Usage: "The cluster to update the services on"},
		&cli.StringFlag{Name: "service", Usage: "Optional service name to update"},
		&cli.BoolFlag{Name: "dry-run", Value: false, Usage: "Parse the template without running upload"},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	sess := session.Must(session.NewSession())
	svc := ecs.New(sess)
	cluster := c.String("cluster")
	filename := c.String("file")
	dryrun := c.Bool("dry-run")

	if dryrun {
		_, err := parser.Parse(filename, os.Environ())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse task definition: %s\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Template %s is parsed successfully\n", filename)
		os.Exit(0)
	}

	taskDefinition := uploadTask(svc, filename)

	if serviceName := c.String("service"); serviceName != "" {
		updateService(svc, serviceName, cluster, taskDefinition)
	}

	return nil
}

func uploadTask(svc *ecs.ECS, filename string) string {
	taskDefinitionInput, err := parser.Parse(filename, os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse task definition: %s\n", err)
		os.Exit(1)
	}

	log.Printf("Registering a task for %s\n", *taskDefinitionInput.Family)
	resp, err := svc.RegisterTaskDefinition(taskDefinitionInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to register task definition: %s\n", err)
		os.Exit(1)
	}

	log.Printf("Created %s\n", *resp.TaskDefinition.TaskDefinitionArn)

	return *resp.TaskDefinition.TaskDefinitionArn
}

func updateService(svc *ecs.ECS, service, cluster, taskDefinition string) {
	log.Printf("Updating service %s\n", service)

	_, err := svc.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        &cluster,
		Service:        &service,
		TaskDefinition: &taskDefinition,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update service %s on cluster %s: %s\n", service, cluster, err)
		os.Exit(1)
	}

	pollUntilTaskDeployed(svc, cluster, service, taskDefinition)
}

func getService(svc *ecs.ECS, service, cluster string) (*ecs.Service, error) {
	resp, err := svc.DescribeServices(&ecs.DescribeServicesInput{
		Services: []*string{aws.String(service)},
		Cluster:  aws.String(cluster),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Failures) > 0 {
		return nil, errors.New(*resp.Failures[0].Reason)
	}

	if len(resp.Services) != 1 {
		return nil, errors.New("multiple services with the same name")
	}

	return resp.Services[0], nil
}

func pollUntilTaskDeployed(svc *ecs.ECS, service string, cluster string, task string) {
	lastSeen := time.Now().Add(-1 * time.Minute)

	for {
		service, err := getService(svc, cluster, service)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get service: %s\n", err)
			os.Exit(1)
		}

		for i := len(service.Events) - 1; i >= 0; i-- {
			event := service.Events[i]
			if event.CreatedAt.After(lastSeen) {
				log.Println(*event.Message)
				lastSeen = *event.CreatedAt
			}
		}

		if len(service.Deployments) == 1 && *service.Deployments[0].TaskDefinition == task {
			return
		}

		time.Sleep(ECS_POLL_INTERVAL)
	}
}
