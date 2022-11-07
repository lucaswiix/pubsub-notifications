package orm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/lucaswiix/meli/notifications/dto"
)

type DynamoDBORM interface {
	CreateTable(tableName string) error
}

type implDynamoORM struct {
	db        *dynamodb.DynamoDB
	userTable string
}

func NewDynamoDBORM(DB *dynamodb.DynamoDB, userTable string) DynamoDBORM {
	return &implDynamoORM{DB, userTable}
}

// CreateTable creates a table
func (r *implDynamoORM) CreateTable(tableName string) error {
	_, err := r.db.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   &tableName,
	})

	return err
}

// PutItem inserts the struct Person
func (r *implDynamoORM) PutUser(user *dto.User) error {
	_, err := r.db.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(user.ID),
			},
			"IsOptOut": {
				BOOL: aws.Bool(user.IsOptOut),
			},
			"CreatedAt": {
				S: aws.String(user.CreatedAt),
			},
		},
		TableName: &r.userTable,
	})

	return err
}
