// +build example

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"time"
)

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	val_time := time.Now().UnixNano() / 1000000

	cfg := Config{}
	if err := cfg.Load(); err != nil {
		exitWithError(fmt.Errorf("failed to load config, %v", err))
	}

	// Create the config specifiing the Region for the DynamoDB table.
	// If Config.Region is not set the region must come from the shared
	// config or AWS_REGION environment variable.
	awscfg := &aws.Config{}
	if len(cfg.Region) > 0 {
		awscfg.WithRegion(cfg.Region)
	}

	// Create the session that the DynamoDB service will use.
	sess, err := session.NewSession(awscfg)
	if err != nil {
		fmt.Println("failed to create session, %v", err)
	} else {
		//	fmt.Println("session success is: ", sess)
	}

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(sess)
	create_params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{ // Required
			{ // Required
				AttributeName: aws.String("camp_devid"), // Required
				AttributeType: aws.String("S"),          // Required
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{ // Required
			{ // Required
				AttributeName: aws.String("camp_devid"), // Required
				KeyType:       aws.String("HASH"),       // Required
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // Required
			ReadCapacityUnits:  aws.Int64(1), // Required
			WriteCapacityUnits: aws.Int64(1), // Required
		},
		TableName: aws.String("dspImpTest1121"), // Required
	}
	resp, err := svc.CreateTable(create_params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("resp is:", resp)
	last_time := time.Now().UnixNano()/1000000 - val_time
	fmt.Println("getItem use  time", last_time)

	//dspImpTest1121

}

type Item struct {
	UserID     int       // Hash key, a.k.a. partition key
	Time       time.Time // Range key, a.k.a. sort ke
	Msg        string    `dynamo:"Message"`
	Count      int       `dynamo:",omitempty"`
	SecretKey  string    `dynamo:"-"` // Ignored
	DeviceId   string    `dynamo:"deviceid"`
	Play       int       `dynamo:"play"`
	CampaginId string    `dynamo:"campid"`
}

type Config struct {
	Table  string // required
	Region string // optional
	Limit  int64  // optional

}

func (c *Config) Load() error {
	//flag.Int64Var(&c.Limit, "limit", 0, "Limit is the max items to be returned, 0 is no limit")
	//flag.StringVar(&c.Table, "table", "", "Table to Query on")
	//flag.StringVar(&c.Region, "region", "", "AWS Region the table is in")
	//flag.Parse()
	c.Limit = 100
	c.Region = "ap-southeast-1"
	//c.Table = "dsp_realimp_test"
	if len(c.Table) == 0 {
		fmt.Println("config is :", c)
		//	flag.PrintDefaults()
		//return fmt.Errorf("table name is required.")
	}

	return nil
}
