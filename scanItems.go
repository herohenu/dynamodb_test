// +build example

package main

import (
	//	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		exitWithError(fmt.Errorf("failed to create session, %v", err))
	}

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(sess)

	// Build the query input parameters

	params := &dynamodb.ScanInput{
		TableName: aws.String(cfg.Table),
	}
	if cfg.Limit > 0 {
		params.Limit = aws.Int64(cfg.Limit)
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		exitWithError(fmt.Errorf("failed to make Query API call, %v", err))
	}
	items := []Item{}
	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	//fmt.Println("resut::::\n", result.Items)
	if err != nil {
		fmt.Println("err is:", err)
	} else {
		fmt.Println("len is: ", len(items))
		//fmt.Println("items is: \n  ", result.Items)
	}

	// Print out the items returned
	for i, item := range items {
		//fmt.Printf("%d:  UserID: %d, Play: %d, Msg: %s  Count: %d SecretKey:%s DeviceId: %s CampaginId:%s \n", i, item.UserID, item.Play, item.Msg, item.Count, item.SecretKey, item.DeviceId, item.CampaginId)
		fmt.Printf("%d:DeviceId: %s Impnum: %d  \n", i, item.DeviceId, item.ImpNum)
	}

	/*
			up_params := &dynamodb.UpdateItemInput{
			Key: map[string]*dynamodb.AttributeValue{ // Required
				"deviceid": { // Required
					S: aws.String("13891389"),
				},
			},
			TableName: aws.String(cfg.Table),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":incr": { // Required
					N: aws.String("1"),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"), //UPDATED_NEW
			UpdateExpression: aws.String("  ADD play :incr"),
		}

		up_resp, up_err := svc.UpdateItem(up_params)

		if up_err != nil {
			fmt.Println("=====>\n update_err:\n", up_err.Error())
		} else {
			fmt.Println("up_res  :========>>>\n ", up_resp)
		}
	*/
	last_time := time.Now().UnixNano()/1000000 - val_time
	fmt.Println(" last time", last_time)

	//
	params_del := &dynamodb.DeleteItemInput{
		TableName: aws.String(cfg.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"deviceid": {
				S: aws.String("111007824_d549646b-96d1-4187-a9e0-662f6d535525"),
			},
		},
	}
	resp_del, err_del := svc.DeleteItem(params_del)
	fmt.Println("err_del i:", err_del)
	fmt.Println("resp_del i:", resp_del)

	//query oneItem  demo
	params_get := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"deviceid": {
				S: aws.String("xxx1231241"),
			},
		},
		TableName: aws.String(cfg.Table), // Required
	}

	resp, err_get := svc.GetItem(params_get)
	oneItem := Item{}
	if err_get == nil {
		// resp is now filled
		fmt.Printf("resp type is %T \n ", resp.Item)
		err = dynamodbattribute.UnmarshalMap(resp.Item, &oneItem)
		if err == nil {
			fmt.Printf("  UserID: %d, Time: %s Msg: %s  Count: %d SecretKey:%s DeviceId: %s CampaginId:%s \n", oneItem.UserID, oneItem.Time, oneItem.Msg, oneItem.Count, oneItem.SecretKey, oneItem.DeviceId, oneItem.CampaginId)
		} else {
			fmt.Println(" Unmarshal err :", err)
		}
		//fmt.Println("convert to Struct obj  err is ", err, "oneItem is:", oneItem)
	} else {
		fmt.Println("GetItem err is: ", err_get)
	}
	last_time = time.Now().UnixNano()/1000000 - val_time
	fmt.Println("getItem use  time", last_time)
	//putItem
	//item := Item {UserID: 6145, Msg: "putMsg", DeviceId: "111-11",CampaginId: "110"}
	params_put := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"deviceid": {
				S: aws.String("key_0013"),
			},
			"UserID": {
				N: aws.String("100100"),
			},
			"Count": {
				N: aws.String("10"),
			},
			"SecretKey": {
				S: aws.String("sk_1112"),
			},
			"play": {
				N: aws.String("10"),
			},
			"userid": {
				N: aws.String("22201"),
			},
			"Message": {
				S: aws.String("put_msg2"),
			},
		},
		TableName: aws.String(cfg.Table), // Required
	}

	resp_put, err_put := svc.PutItem(params_put)
	if err_put != nil {
		fmt.Println("err_put is: ", err_put)
	} else {
		fmt.Println("resp_put is:", resp_put)
	}

	last_time = time.Now().UnixNano()/1000000 - val_time
	fmt.Println("putItem use  time", last_time)
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
	ImpNum     int       `dynamo:"ImpNum,omitempty"`
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
	c.Table = "dsp_realimp_test"
	if len(c.Table) == 0 {
		//	flag.PrintDefaults()
		return fmt.Errorf("table name is required.")
	}

	return nil
}
