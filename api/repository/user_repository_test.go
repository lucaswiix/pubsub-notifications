package repository

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gusaul/go-dynamock"
	"github.com/lucaswiix/meli/notifications/dto"
	"github.com/stretchr/testify/assert"
)

var mock *dynamock.DynaMock

var (
	mockUserID = "34065a27-cb0c-42f9-9bea-258c5806aaa5"
	ctx        = context.TODO()
)

func TestPutUserSuccess(t *testing.T) {
	db, mock := dynamock.New()

	repo := NewUserRepository(db)

	user := dto.User{
		ID:        mockUserID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	insertUser := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(user.ID),
		},
		"IsOptOut": {
			BOOL: aws.Bool(user.IsOptOut),
		},
		"CreatedAt": {
			S: aws.String(user.CreatedAt),
		},
	}

	mock.ExpectPutItem().ToTable("users").WithItems(insertUser).WillReturns(dynamodb.PutItemOutput{})

	err := repo.Put(user, ctx)

	assert.NoError(t, err)
}

func TestDelUserSuccess(t *testing.T) {
	db, mock := dynamock.New()

	repo := NewUserRepository(db)

	user := dto.User{
		ID:        mockUserID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	expectDel := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(user.ID),
		},
	}

	mock.ExpectDeleteItem().ToTable("users").WithKeys(expectDel).WillReturns(dynamodb.DeleteItemOutput{})

	err := repo.Del(user.ID, ctx)

	assert.NoError(t, err)
}

func TestGetUserSuccess(t *testing.T) {
	db, mock := dynamock.New()

	repo := NewUserRepository(db)

	userMock := dto.User{
		ID:        mockUserID,
		IsOptOut:  true,
		CreatedAt: time.Now().Format(time.RFC822),
	}

	expectGet := map[string]*dynamodb.AttributeValue{
		"Id": {
			S: aws.String(userMock.ID),
		},
	}

	result := dynamodb.GetItemOutput{

		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(userMock.ID),
			},
			"IsOptOut": {
				BOOL: aws.Bool(userMock.IsOptOut),
			},
			"CreatedAt": {
				S: aws.String(userMock.CreatedAt),
			},
		},
	}

	mock.ExpectGetItem().ToTable("users").WithKeys(expectGet).WillReturns(result)

	user, err := repo.Get(userMock.ID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, user, userMock)
}
