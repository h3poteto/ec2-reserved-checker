# ec2-reserved-checker

`ec2-reserved-checker` is a management tool for AWS EC2 Reserved Instances. It show running EC2 Instances which is not applied Resereved Instance, and it show active Reserved Instances which does not relate any running EC2 Instances.

## Install
Get binary from github:

```
$ wget https://github.com/h3poteto/ec2-reserved-checker/releases/download/v0.1.0/ec2_reserved_checker_0.1.0_darwin_adm64.zip
```

or, build. It requires Go1.6 and [gom](https://github.com/mattn/gom).

```
$ git clone git@github.com:h3poteto/ec2-reserved-checker.git
$ cd ec2-reserved_checker
$ export GO15VENDOREXPERIMENT=0  # It is required in gom, when we use go1.6: https://github.com/mattn/gom/issues/80
$ gom install
$ gom build -o ec2-reserved-checker main.go
```

## Setup
It use [aws-sdk-go](https://github.com/aws/aws-sdk-go), so please set enviroments for AWS:

```
$ export AWS_ACCESS_KEY_ID=AKID1234567890
$ export AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

or, ensure that you've configured credentials in `~/.aws/credentials` :

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

And, `ec2-reserved-checker` requires `AWS_REGION` to search EC2 Instances in your AWS Account.
```
$ export AWS_REGION=ap-northeast-1
```

### Required IAM permissions
The user/role used should have `ec2:DescribeInstances` and `ec2:DescribeReservedInstances`
permissions. So you should configure one policy like:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:DescribeReservedInstances"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## Example

```
$ ./ec2-reserved-checker
----------------------------------------------
 There are 4 running EC2 instances
----------------------------------------------
  EC2 Instance ID: i-0a8e47ff, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro
  EC2 Instance ID: i-0c66cd221e4c87f20, AvailabilityZone: ap-northeast-1c, InstanceType: t2.micro
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

### Options 
```
$ ./ec2-reserved-checker -h
Usage of ./ec2-reserved-checker:
  -info
    	Show all instances information (default true)
  -notapplied
    	Show not applied ec2 instances information (default true)
  -unused
    	Show unused reserved instances information (default true)

$ ./ec2-reserved-checker -info=false -unused=false
----------------------------------------------
 There are 1 EC2 Instances which are not applied reserved
----------------------------------------------
  EC2 Instance ID: i-fd5e1a58, AvailabilityZone: ap-northeast-1a, InstanceType: t2.micro

```


## Author

[Akira Fukushima](https://github.com/h3poteto)

