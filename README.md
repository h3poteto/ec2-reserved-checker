# ec2-reserved-checker

For example: 

```
$ gom run main.go
----------------------------------------------
 There are 4 running EC2 instances
----------------------------------------------
  EC2 Instance ID: i-0a8e47ff, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro
  EC2 Instance ID: i-d6811259, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro
  EC2 Instance ID: i-be8a454b, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro
  EC2 Instance ID: i-fd5e1a58, AvailabilityZone: ap-northeast-1a, InstanceType: t2.micro

----------------------------------------------
There are 2 active Reserved instances
----------------------------------------------
  Reserved Instance ID: c4d57437-4c84-4e6b-9d77-ed9d79926fe3, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro, number: 2
  Reserved Instance ID: d02bab08-6e46-4627-ad46-530db72207de, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro, number: 1

----------------------------------------------
 There are 1 EC2 Instances which are not applied reserved
----------------------------------------------
  EC2 Instance ID: i-fd5e1a58, AvailabilityZone: ap-northeast-1a, InstanceType: t2.micro

----------------------------------------------
 There are 0 Reserved Instances which are not related running EC2 instances
----------------------------------------------


```
