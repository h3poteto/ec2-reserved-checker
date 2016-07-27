package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	instances, reservedInstances, err := EC2InstancesAndReservedInstances()
	if err != nil {
		panic(err)
	}

	fmt.Println("----------------------------------------------")
	fmt.Printf(" There are %v running EC2 instances\n", len(instances))
	fmt.Println("----------------------------------------------")

	for _, inst := range instances {
		fmt.Printf("EC2 Instance ID: %v\n", *inst.InstanceId)
	}

	fmt.Println("----------------------------------------------")
	fmt.Printf("There are %v active Reserved instances\n", len(reservedInstances))
	fmt.Println("----------------------------------------------")
	for _, inst := range reservedInstances {
		fmt.Printf("Reserved Instance ID: %v, number: %v\n", *inst.ReservedInstancesId, *inst.InstanceCount)
	}

	_ = flattenReservedInstances(reservedInstances)

	// reservedAttachedInstances := make([]*ec2.Instance, len(instances))
	// copy(unusedInstances, instances)
	// usedReservedInstances := make([]*ec2.ReservedInstances, 0)

	// // Need EC2 Instances which not related any active Reserved Instances
	// //
	// for _, ri := range reservedInstances {
	// 	for i := 0; i < int(*ri.InstanceCount); i++ {
	// 		for j, inst := range instances {
	// 			if *ri.AvailabilityZone == *inst.Placement.AvailabilityZone && *ri.InstanceType == *inst.InstanceType {
	// 				if j < (len(instances) - 1) {
	// 					unusedInstances = append(unusedInstances[:j], unusedInstances[(j+1):]...)
	// 				} else {
	// 					unusedInstances = unusedInstances[:j]
	// 				}
	// 				usedReservedInstances = append(usedReservedInstances, ri)
	// 				break
	// 			}
	// 		}
	// 	}
	// }
}

// EC2InstancesAndReservedInstances get running EC2 Instances and active Reserved Instances
func EC2InstancesAndReservedInstances() ([]*ec2.Instance, []*ec2.ReservedInstances, error) {
	region := os.Getenv("AWS_REGION")
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})

	// Get running instances
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
		return nil, nil, err
	}

	var runningEC2Instances []*ec2.Instance

	// instanceResp has all of the response data, pull out instance IDs:
	for idx, _ := range instanceResp.Reservations {
		for _, inst := range instanceResp.Reservations[idx].Instances {
			runningEC2Instances = append(runningEC2Instances, inst)
		}
	}

	// Get active reserved instances
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

	// Call the DescribeReservedInstances Operation
	reservedResp, err := svc.DescribeReservedInstances(reservedParams)
	if err != nil {
		return nil, nil, err
	}

	activeReservedInstances := reservedResp.ReservedInstances

	return runningEC2Instances, activeReservedInstances, nil
}

func flattenReservedInstances(reservedInstances []*ec2.ReservedInstances) []*ec2.ReservedInstances {
	ri := make([]*ec2.ReservedInstances, 0)
	for _, inst := range reservedInstances {
		for i := 0; i < int(*inst.InstanceCount); i++ {
			ri = append(ri, inst)
		}
	}
	return ri
}
