# example fargate ecs service


### networking
Create a VPC and some subnets, or use the default VPC. This is a complicated topic covered by the [aws docs](http://docs.aws.amazon.com/AmazonVPC/latest/GettingStartedGuide/getting-started-ipv4.html)

All of the following templates will need the vpc and subnets replaced.

Additionally, they also need two security groups pre created.

A security group for the ALB and a security group for the fagate containers. The ELBs should be able to access container ports. Update the security group references in the templates.

### upload [fargate-task-definition.yml](fargate-task-definition.yml) 

This will need to be modified to fit your use case. Most important is the task role. This should provide access to pull containers from ECR as well as write cloudwatch logs, as well as whatever AWS perms your app needs.

```bash
ecs-upload-task --file fargate-task-definition.yml
```

### create the ELB

```bash
aws cloudformation create-stack \
    --stack-name myApp-alb \
    --template-body file://alb.yml
```

### create the ECS service

```
aws cloudformation create-stack \
    --stack-name myApp \
    --template-body file://fargate-task-definition.yml \
    --parameters '[{"ParameterKey":"TaskDefinition","ParameterValue":"arn:aws:ecs:us-east-1:447214301260:task-definition/myapp:{OUTPUT FROM ecs-upload-task}"}]'
```

At this point everything should be running. 

### deploying

Later on in your CI pipeline, after pushing a new container up:
```bash
IMAGE_NAME=your/service:r1234 ecs-upload-task \
    --file fargate-task-definition.yml
    --service myApp-prod-1234567 # you can find this in the ECS web ui
```