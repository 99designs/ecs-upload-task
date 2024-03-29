# `ecs-upload-task`

A very simple cli tool that will upload an ecs task and optionally update a service to reference it.

Designed as a partner tool to https://github.com/buildkite/ecs-run-task and shares the same task definition parser with environment substitution.

## Example

Install it
```bash
go install github.com/99designs/ecs-upload-task@latest

or

curl -L --fail -o /usr/bin/ecs-upload-task \
  <release from https://github.com/99designs/ecs-upload-task/releases>
```

`taskdefinition.yml`
```yaml
family: hello-world
containerDefinitions:
- essential: true
  image: ${IMAGE_NAME:-99designs/hello-world}
  memory: 64
  cpu: 10
  name: hello-world
  portMappings:
    - {containerPort: 1234, hostPort: 1234}
  environment:
    - {name: PORT, value: "1234"}

```

then take it for a dry-run to validate that the template is parsed successfully.

```bash
ecs-upload-task --file taskdefinition.yml --dry-run
Template taskdefinition.yml is parsed successfully
```

then register the task definition

```bash
ecs-upload-task --file taskdefinition.yml
```

Create a new ECS service referencing the task definition. Service templates can use the family name without the version, simply `hello-world` is enough.


Later, in your automated CI pipeline you can deploy by simply:

```bash
export IMAGE_NAME=99designs/hello-world:$FRESH_BUILD_NUMBER
ecs-upload-task --file taskdefinition.yml --service hello-world-2017-05-15-10-45

2017/11/14 18:02:32 Registering a task for hello-world
2017/11/14 18:02:34 Created arn:xxx:task-definition/hello-world:2
2017/11/14 18:02:34 Updating service hello-world-2017-05-15-10-45
2017/11/14 18:02:55 (service hello-world-2017-05-15-10-45) has started 1 tasks: (task 81b2963f-072a-479b-856f-26af2ec615f8).
2017/11/14 18:03:05 (service hello-world-2017-05-15-10-45) registered 1 instances in (elb hello-world-ELB-O5IUREC150O5)
2017/11/14 18:03:28 (service hello-world-2017-05-15-10-45) deregistered 1 instances in (elb hello-world-ELB-O5IUREC150O5)
2017/11/14 18:03:48 (service hello-world-2017-05-15-10-45) has stopped 1 running tasks: (task a3e3a91b-be05-4092-bcf6-47f2075933af).

```

## Creating a new release

1. Check out the the commit you want to create a release for, and tag it with appropriate semver convention:

```
$ git tag vx.x.x
$ git push --tags
```

2. Create the binaries:

Binaries are compiled using [`gox`](https://github.com/mitchellh/gox). This provides a simple way to compile for multiple CPU architectures (`amd64`, `arm64`), which allows us to provide support for devices such as M1 Macs.

`gox` can be installed via `go install`

```
go install github.com/mitchellh/gox@latest
```
Then simply run the `Makefile` to compile.

```
$ make
```

3. Go to https://github.com/99designs/ecs-upload-task/releases/new

4. Select the tag version you just created

5. Attach the binaries from `./bin/*`

