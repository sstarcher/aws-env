package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func shellEncode(key, value string) {
	key = strings.ToUpper(key)
	key = strings.Replace(key, ":", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	fmt.Printf("AWS_%s=%s\n", strings.Replace(strings.ToUpper(key), ":", "_", -1), value)
}

func main() {
	session := session.New()
	metadata := ec2metadata.New(session)

	instanceID, err := metadata.GetMetadata("instance-id")
	if err != nil {
		panic(err)
	}
	shellEncode("INSTANCE_ID", instanceID)

	region, err := metadata.Region()
	if err != nil {
		panic(err)
	}
	shellEncode("REGION", region)

	ec2client := ec2.New(session, aws.NewConfig().WithRegion(region))
	resp, err := ec2client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})
	if err != nil {
		panic(err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			for _, tag := range inst.Tags {
				shellEncode("Tag_"+*tag.Key, *tag.Value)
			}
		}
	}
}
