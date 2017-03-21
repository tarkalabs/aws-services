package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/tarkalabs/aws-services/models"
	"log"
	"sync"
)

var conn *sqlx.DB

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func getInstances(reg *ec2.Region, wg *sync.WaitGroup) {
	fmt.Println("Displaying instances for ", *reg.RegionName)
	sess, err := session.NewSession()
	failOnError(err, "unable to create session")
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
	wg.Done()
}
func main() {
	// initialize database
	done := make(chan bool)
	conn = models.InitDb()

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	regionsSvc := ec2.New(sess, aws.NewConfig().WithRegion("us-east-1"))
	regions, err := regionsSvc.DescribeRegions(nil)
	if err != nil {
		fmt.Println("Unable to list regions", err)
	}
	var wg sync.WaitGroup
	for _, reg := range regions.Regions {
		wg.Add(1)
		go getInstances(reg, &wg)
	}

	go func() {
		wg.Wait()
		s3Svc := s3.New(sess, aws.NewConfig().WithRegion("us-east-1"))
		bucketList, err := s3Svc.ListBuckets(nil)
		failOnError(err, "Error listing Buckets")
		for _, bucket := range bucketList.Buckets {
			fmt.Println(*bucket.Name, bucket)

			s3_attrs, err := json.Marshal(bucket)
			failOnError(err, "Error getting Bucket attributes")

			fmt.Println(string(s3_attrs))

		}
		done <- true

	}()
	<-done
}
