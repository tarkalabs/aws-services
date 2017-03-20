package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/tarkalabs/aws-services/models"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// initialize database
	conn := models.InitDb()
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	regionsSvc := ec2.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	regions, err := regionsSvc.DescribeRegions(nil)
	if err != nil {
		fmt.Println("Unable to list regions", err)
	}
	for _, reg := range regions.Regions {
		fmt.Println("Displaying instances for ", *reg.RegionName)

		svc := ec2.New(sess, aws.NewConfig().WithRegion(*reg.RegionName))
		resp, err := svc.DescribeInstances(nil)

		if err != nil {
			panic(err)
		}
		for _, r := range resp.Reservations {
			fmt.Printf("Found %d instances.\n", len(r.Instances))
			for _, inst := range r.Instances {
				i := models.NewInstance()
				i.Name = *inst.InstanceId
				i.Region = *reg.RegionName
				attrs, err := json.Marshal(inst)
				failOnError(err, "unabled to encode json")
				fmt.Println(string(attrs))
				i.Attributes = string(attrs)
				models.SaveInstance(conn, i)
				fmt.Println("found instance : ", *inst.InstanceId)
			}
		}
	}
}
