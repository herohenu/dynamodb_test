// +build example

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"strconv"
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
	///////////////////////////////
	//deleteTable(svc, "ImpTest1122")

	fmt.Printf(" type is :%T", svc)
	//createTable(svc, "ImpTest1122")
	//	batchPutItems(svc)
	batchGetItems(svc)
	last_time := time.Now().UnixNano()/1000000 - val_time
	fmt.Println("create table use  time", last_time)

}
func batchGetItems(svc *dynamodb.DynamoDB) {
	requestItems := buildRequestItems([]string{"1", "2", "3", "4"})
	params := &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}

	resp, err := svc.BatchGetItem(params)
	if err != nil {
		fmt.Println("batch getItems err :", err)
	}

	fmt.Println("  getTable res::", resp.Responses["ImpTest1122"])
	items := []Item{}
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Responses["ImpTest1122"], &items)
	fmt.Println("batchgetItems  resp is convert toItems is  ::::  ", items)

}
func deleteTable(svc *dynamodb.DynamoDB, tableName string) {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName), // Required
	}
	resp, err := svc.DeleteTable(params)
	fmt.Println("delete resp is: ", resp, "err :", err)

}

func buildRequestItems(ids []string) map[string]*dynamodb.KeysAndAttributes {
	//ids = []string{"1", "2", "3"}
	makeKeysAndAttrs := func() *dynamodb.KeysAndAttributes {
		out := &dynamodb.KeysAndAttributes{Keys: []map[string]*dynamodb.AttributeValue{}}
		for _, r := range ids {
			out.Keys = append(out.Keys, map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(r)}})
		}
		out.AttributesToGet = append(out.AttributesToGet, aws.String("Val"))
		out.AttributesToGet = append(out.AttributesToGet, aws.String("Id"))
		//	fmt.Println(" out is; ", out)
		return out
	}
	//fmt.Println(makeKeysAndAttrs())
	return map[string]*dynamodb.KeysAndAttributes{"ImpTest1122": makeKeysAndAttrs()}
}

/*key is 10 .199*/
func batchPutItems(svc *dynamodb.DynamoDB) {
	for i := 1; i < 10; i++ {
		key := strconv.Itoa(i)
		fmt.Println("key  is: ", key)
		putItem(svc, key)
	}

}
func putItem(svc *dynamodb.DynamoDB, key string) {
	//"campId001_devId111")
	params_put := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id":  {S: aws.String(key)},
			"Val": {S: aws.String("1")},
		},
		TableName: aws.String("ImpTest1122"), // Required
	}
	resp_put, err_put := svc.PutItem(params_put)
	if err_put != nil {
		fmt.Println("err_put is: ", err_put)
	} else {
		fmt.Println("resp_put is:", resp_put)
	}

}

func createTable(svc *dynamodb.DynamoDB, tableName string) {
	create_params := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{ // Required
			{ // Required
				AttributeName: aws.String("Id"), // Required
				AttributeType: aws.String("S"),  // Required
			},
			/*{ // Required
				AttributeName: aws.String("Val"), // Required
				AttributeType: aws.String("S"),   // Required
			},*/
		},
		KeySchema: []*dynamodb.KeySchemaElement{ // Required
			{ // Required
				AttributeName: aws.String("Id"),   // Required
				KeyType:       aws.String("HASH"), // Required
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ // Required
			ReadCapacityUnits:  aws.Int64(1), // Required
			WriteCapacityUnits: aws.Int64(1), // Required
		},
		TableName: aws.String(tableName), // Required
	}
	resp, err := svc.CreateTable(create_params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("resp is:", resp)

}
func updateItem(svc *dynamodb.DynamoDB) {
	up_params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"Id": {S: aws.String("campId001_devId111")},
		},
		TableName: aws.String("ImpTest1122"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": { // Required
				N: aws.String("1"),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"), //UPDATED_NEW
		UpdateExpression: aws.String("  ADD getN(Val) :incr"),
	}

	up_resp, up_err := svc.UpdateItem(up_params)

	if up_err != nil {
		fmt.Println("=====>\n update_err:\n", up_err.Error())
	} else {
		fmt.Println("up_res  :========>>>\n ", up_resp)
	}
}

type Item struct {
	Id  string
	Val string
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
	c.Table = "ImpTest1122"
	if len(c.Table) == 0 {
		fmt.Println("config is :", c)
		//	flag.PrintDefaults()
		//return fmt.Errorf("table name is required.")
	}

	return nil
}
