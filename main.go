package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type FlattenReservedInstances struct {
	Reserved *ec2.ReservedInstances
	Index    int
}

func main() {

	runningInstances, reservedInstances, err := EC2InstancesAndReservedInstances()
	if err != nil {
		panic(err)
	}

	fmt.Println("----------------------------------------------")
	fmt.Printf(" There are %v running EC2 instances\n", len(runningInstances))
	fmt.Println("----------------------------------------------")

	for _, instance := range runningInstances {
		fmt.Printf("  EC2 Instance ID: %v, AvailabilityZone: %v, InstanceType: %v\n", *instance.InstanceId, *instance.Placement.AvailabilityZone, *instance.InstanceType)
	}
	fmt.Println("")

	fmt.Println("----------------------------------------------")
	fmt.Printf("There are %v active Reserved instances\n", len(reservedInstances))
	fmt.Println("----------------------------------------------")
	for _, reserved := range reservedInstances {
		fmt.Printf("  Reserved Instance ID: %v, AvailabilityZone: %v, InstanceType: %v, number: %v\n", *reserved.ReservedInstancesId, *reserved.AvailabilityZone, *reserved.InstanceType, *reserved.InstanceCount)
	}
	fmt.Println("")

	reservedAppliedInstances := make([]*ec2.Instance, 0)
	relatedReservedInstances := make([]*FlattenReservedInstances, 0)

	flattenReservedInstances := flattenReservedInstances(reservedInstances)

	// reserved instances which related to EC2 instances
	for _, flattenReserved := range flattenReservedInstances {
		for _, instance := range runningInstances {
			if *flattenReserved.Reserved.AvailabilityZone == *instance.Placement.AvailabilityZone && *flattenReserved.Reserved.InstanceType == *instance.InstanceType {
				reservedAppliedInstances = append(reservedAppliedInstances, instance)
				relatedReservedInstances = append(relatedReservedInstances, flattenReserved)
			}
		}
	}

	// We need ondemand instances which are not applied reserved,
	// and reserved instances which are not related running EC2 instances.
	reservedNotAppliedInstances := make([]*ec2.Instance, 0)
	unusedReservedInstances := make([]*FlattenReservedInstances, 0)

	for _, instance := range runningInstances {
		applied := false
		for _, appliedInst := range reservedAppliedInstances {
			if *instance.InstanceId == *appliedInst.InstanceId {
				applied = true
				break
			}
		}
		if !applied {
			reservedNotAppliedInstances = append(reservedNotAppliedInstances, instance)
		}
	}

	for _, flattenReserved := range flattenReservedInstances {
		related := false
		for _, relatedReserved := range relatedReservedInstances {
			if flattenReserved.Index == relatedReserved.Index {
				related = true
				break
			}
		}
		if !related {
			unusedReservedInstances = append(unusedReservedInstances, flattenReserved)
		}
	}

	fmt.Println("----------------------------------------------")
	fmt.Printf(" There are %v EC2 Instances which are not applied reserved\n", len(reservedNotAppliedInstances))
	fmt.Println("----------------------------------------------")
	for _, instance := range reservedNotAppliedInstances {
		fmt.Printf("  EC2 Instance ID: %v, AvailabilityZone: %v, InstanceType: %v\n", *instance.InstanceId, *instance.Placement.AvailabilityZone, *instance.InstanceType)
	}
	fmt.Println("")

	fmt.Println("----------------------------------------------")
	fmt.Printf(" There are %v Reserved Instances which are not related running EC2 instances\n", len(unusedReservedInstances))
	fmt.Println("----------------------------------------------")
	for _, flattenReserved := range unusedReservedInstances {
		fmt.Printf("  Reserved Instance ID: %v, AvailabilityZone: %v, InstanceType: %v, number: 1\n", *flattenReserved.Reserved.ReservedInstancesId, *flattenReserved.Reserved.AvailabilityZone, *flattenReserved.Reserved.InstanceType)
	}
	fmt.Println("")
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

func flattenReservedInstances(reservedInstances []*ec2.ReservedInstances) []*FlattenReservedInstances {
	flattenReserved := make([]*FlattenReservedInstances, 0)
	duplicateReservedInstances := make([]*ec2.ReservedInstances, 0)
	for _, reserved := range reservedInstances {
		for i := 0; i < int(*reserved.InstanceCount); i++ {
			duplicateReservedInstances = append(duplicateReservedInstances, reserved)
		}
	}

	for i, reserved := range duplicateReservedInstances {
		flatten := &FlattenReservedInstances{
			Reserved: reserved,
			Index:    i,
		}
		flattenReserved = append(flattenReserved, flatten)
	}

	return flattenReserved
}
