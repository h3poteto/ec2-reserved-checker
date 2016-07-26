package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	region := os.Getenv("AWS_REGION")
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	instancesParams := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	// Call the DescribeInstances Operation
	instanceResp, err := svc.DescribeInstances(instancesParams)
	if err != nil {
		panic(err)
	}

	var allEC2Instances []*ec2.Instance

	// instanceResp has all of the response data, pull out instance IDs:
	fmt.Println("> Number of reservation sets: ", len(instanceResp.Reservations))
	for idx, _ := range instanceResp.Reservations {
		for _, inst := range instanceResp.Reservations[idx].Instances {
			allEC2Instances = append(allEC2Instances, inst)
		}
	}

	for _, instance := range allEC2Instances {
		fmt.Printf(" Instance ID: %v, state: %v\n", *instance.InstanceId, *instance.State.Name)
	}

	reservedParams := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("state"),
				Values: []*string{
					aws.String("active"),
				},
			},
		},
	}

	reservedResp, err := svc.DescribeReservedInstances(reservedParams)
	if err != nil {
		panic(err)
	}

	allReservedInstances := reservedResp.ReservedInstances

	for _, instance := range allReservedInstances {
		fmt.Printf(" ReservedInstances ID: %v\n", *instance.ReservedInstancesId)
	}

}
