package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	flag "github.com/ogier/pflag"
)

// findClusterArn finds a cluster's arn from a more human
// readable name
func findClusterArn(svc *ecs.ECS, clusterName string) (string, error) {
	params := &ecs.ListClustersInput{}
	clusters, err := svc.ListClusters(params)
	if err != nil {
		return "", err
	}
	for _, ClusterID := range clusters.ClusterArns {
		var pattern = `/` + clusterName + `$`
		matched, _ := regexp.MatchString(pattern, *ClusterID)
		if matched {
			return *ClusterID, nil
		}
	}
	return "", errors.New("Could not find cluster with name: " + clusterName)
}

// findServiceArn finds a service arn from a more human
// readable name
func findServiceArn(svc *ecs.ECS, clusterArn string, serviceName string) (string, error) {
	params := &ecs.ListServicesInput{
		Cluster: aws.String(clusterArn),
	}
	params.SetMaxResults(100)
	services, err := svc.ListServices(params)
	if err != nil {
		return "", err
	}
	for _, ServiceID := range services.ServiceArns {
		var pattern = `/` + serviceName + `$`
		matched, _ := regexp.MatchString(pattern, *ServiceID)
		if matched {
			return *ServiceID, nil
		}
	}
	return "", errors.New("Could not find service with name: " + serviceName)
}

// findTaskArn
// finds the first task definition in a cluster and service.
//
// Returns task ARN
//
// TODO:
//   filter the resulting tasks by essential: true
//   and only return the essential?
//   either that or assume that the new image will
//   somewhat match the old image, and use that to pick
//   one of the tasks.
func findTaskArn(svc *ecs.ECS, clusterArn string, serviceArn string) (string, error) {
	params := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(serviceArn),
		},
		Cluster: aws.String(clusterArn),
	}
	resp, err := svc.DescribeServices(params)
	return *resp.Services[0].TaskDefinition, err
}

// setImage
// Create a new Task Definition based on an existing
// ARN, and a new image.
//
// Returns new task's ARN
//
func setImage(svc *ecs.ECS, taskArn string, image string) (string, error) {
	params := &ecs.DescribeTaskDefinitionInput{TaskDefinition: aws.String(taskArn)}
	resp, err := svc.DescribeTaskDefinition(params)
	if err != nil {
		return "", err
	}
	task := resp.TaskDefinition
	var out = taskArn
	if *task.ContainerDefinitions[0].Image != image {
		task.ContainerDefinitions[0].Image = &image
		regResp, err := svc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
			Family:               task.Family,
			ContainerDefinitions: task.ContainerDefinitions,
			Volumes:              task.Volumes,
		})
		if err != nil {
			return "", err
		}
		out = *regResp.TaskDefinition.TaskDefinitionArn
	}
	return out, nil
}

func main() {

	var cluster = flag.String("cluster", "default", "Name of cluster")

	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage:\n\taws-rollout [service] [image]")
		return
	}

	var service = flag.Arg(0)
	var image = flag.Arg(1)

	fmt.Printf("Cluster: %s\n", *cluster)
	fmt.Printf("Service: %s\n", service)
	fmt.Printf("Image: %s\n", image)

	svc := ecs.New(session.New())

	clusterArn, err := findClusterArn(svc, *cluster)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	serviceArn, err := findServiceArn(svc, clusterArn, service)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	taskArn, err := findTaskArn(svc, clusterArn, serviceArn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	newTaskArn, err := setImage(svc, taskArn, image)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	params := &ecs.UpdateServiceInput{
		Service:        aws.String(serviceArn),
		Cluster:        aws.String(clusterArn),
		TaskDefinition: aws.String(newTaskArn),
	}
	serv, err := svc.UpdateService(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Deployed Task: %s\n", newTaskArn)
	fmt.Printf("Pending Count: %d\n", *serv.Service.PendingCount)
	fmt.Println("Deployment Success!")
}
