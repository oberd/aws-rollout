package main

import (
    "fmt"
    "errors"
    "regexp"
    flag "github.com/ogier/pflag"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ecs"
)

// findClusterArn finds a cluster's arn from a more human
// readable name
func findClusterArn(svc *ecs.ECS, clusterName string) (string, error)  {
    params := &ecs.ListClustersInput{}
    clusters, err := svc.ListClusters(params)
    if err != nil {
        return "", err
    }
    for _, clusterId := range clusters.ClusterArns {
        var pattern string = `/` + clusterName + `$`;
        matched, _ := regexp.MatchString(pattern, *clusterId)
        if matched {
            return *clusterId, nil
        }
    }
    return "", errors.New("Could not find cluster with name: " + clusterName)
}

// findServiceArn finds a service arn from a more human
// readable name
func findServiceArn(svc *ecs.ECS, clusterArn string, serviceName string) (string, error)  {
    params := &ecs.ListServicesInput{
        Cluster: aws.String(clusterArn),
    }
    services, err := svc.ListServices(params)
    if err != nil {
        return "", err
    }
    for _, serviceId := range services.ServiceArns {
        var pattern string = `/` + serviceName + `$`;
        matched, _ := regexp.MatchString(pattern, *serviceId)
        if matched {
            return *serviceId, nil
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
    params := &ecs.DescribeTaskDefinitionInput{ TaskDefinition: aws.String(taskArn) }
    resp, err := svc.DescribeTaskDefinition(params)
    if err != nil {
        return "", err
    }
    task := resp.TaskDefinition
    task.ContainerDefinitions[0].Image = &image
    regResp, err := svc.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
        Family: task.Family,
        ContainerDefinitions: task.ContainerDefinitions,
        Volumes: task.Volumes,
    })
    return *regResp.TaskDefinition.TaskDefinitionArn, nil
}

func main() {

    var cluster *string = flag.String("cluster", "default", "Name of cluster")
    
    flag.Parse()

    if flag.NArg() < 2 {
        fmt.Println("Usage:\n\taws-rollout [service] [image]")
        return
    }

    var service string = flag.Arg(0)
    var image string = flag.Arg(1)

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
    newTaskArn, err := setImage(svc, taskArn, image)
    params := &ecs.UpdateServiceInput{
        Service: aws.String(serviceArn),
        Cluster: aws.String(clusterArn),
        TaskDefinition: aws.String(newTaskArn),
    }
    serv, err := svc.UpdateService(params)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Printf("Deployed %s %d", newTaskArn, *serv.Service.PendingCount)
}