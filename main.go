package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func shellEncode(key, value string, output io.Writer) {
	key = strings.ToUpper(key)
	key = strings.Replace(key, ":", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	fmt.Fprintf(output, "AWS_%s=%s\n", strings.Replace(strings.ToUpper(key), ":", "_", -1), value)
}

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	session := session.New()
	metadata := ec2metadata.New(session)
	f, err := os.Create("/etc/aws")
	if err != nil {
		log.Panicf("Failed to create file /etc/aws, %v", err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)

	var instanceID string
	if metadata.Available() {
		instanceID, err = metadata.GetMetadata("instance-id")
		if err != nil {
			log.Panicf("Failed to fetch instance-id, %v", err)
		}
		shellEncode("INSTANCE_ID", instanceID, writer)

		az, err := metadata.GetMetadata("placement/availability-zone")
		if err != nil {
			log.Panicf("Failed to fetch availablity-zone, %v", err)
		}
		shellEncode("AVAILABLITY_ZONE", az, writer)

		region, err := metadata.Region()
		if err != nil {
			log.Panicf("Failed to fetch Region, %v", err)
		}
		shellEncode("REGION", region, writer)
		session = session.Copy(&aws.Config{Region: aws.String(region)})
	} else {
		instanceID = os.Args[1]
	}

	iamClient := iam.New(session)
	aliases, err := iamClient.ListAccountAliases(&iam.ListAccountAliasesInput{})
	if err != nil {
		log.Panicf("Failed to ListAccountAliases, %v", err)
	}

	for _, alias := range aliases.AccountAliases {
		shellEncode("ACCOUNT_ALIAS", *alias, writer)
	}

	log.Info("Describing instances")
	ec2Client := ec2.New(session)
	resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})
	if err != nil {
		log.Panicf("Failed to DescribeInstances, %v", err)
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			log.Infof("Instance ID: %s", *inst.InstanceId)
			for _, tag := range inst.Tags {
				shellEncode("Tag_"+*tag.Key, *tag.Value, writer)
			}
		}
	}

	writer.Flush()
}
