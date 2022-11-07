package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/lucaswiix/meli/notifications/dto"
	"go.elastic.co/apm"
)

//go:generate mockgen -destination=mock/opt_out.go -package=mock . UserRepository
type UserRepository interface {
	Put(user dto.User, ctx context.Context) error
	Del(userID string, ctx context.Context) error
	Get(userID string, ctx context.Context) (dto.User, error)
}

type implUserRepository struct {
	DB dynamodbiface.DynamoDBAPI
}

func NewUserRepository(DB dynamodbiface.DynamoDBAPI) UserRepository {
	return &implUserRepository{DB}
}

func (r *implUserRepository) Put(user dto.User, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SetUser", "repository")
	defer span.End()

	_, err := r.DB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
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
		TableName: &UsersTable,
	})

	return err
}

func (r *implUserRepository) Del(userID string, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "DelUser", "repository")
	defer span.End()
	_, err := r.DB.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userID),
			},
		},
		TableName: &UsersTable,
	})

	return err
}

func (r *implUserRepository) Get(userID string, ctx context.Context) (dto.User, error) {
	span, ctx := apm.StartSpan(ctx, "IsOptOut", "repository")
	defer span.End()
	result, err := r.DB.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userID),
			},
		},
		TableName: &UsersTable,
	})

	if err != nil {
		return dto.User{}, err
	}
	user := dto.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	return user, err
}
