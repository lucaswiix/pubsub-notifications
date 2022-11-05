package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestOptOutSetSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectSet(fmt.Sprintf("opt-out:%s", userUUID), optOutVal, 0).SetVal("")

	err := repo.Set(userUUID, ctx)

	assert.NoError(t, err)
}

func TestOptOutSetError(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectSet(fmt.Sprintf("opt-out:%s", userUUID), optOutVal, 0).SetErr(errors.New("error"))

	err := repo.Set(userUUID, ctx)

	assert.EqualError(t, err, "error")
}

func TestOptOutDelSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectDel(fmt.Sprintf("opt-out:%s", userUUID)).SetVal(1)

	err := repo.Del(userUUID, ctx)

	assert.NoError(t, err)
}

func TestOptOutIsSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectGet(fmt.Sprintf("opt-out:%s", userUUID)).SetVal("true")

	isOptOut, err := repo.Is(userUUID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, isOptOut, true)
}

func TestOptOutIsFalseSuccess(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectGet(fmt.Sprintf("opt-out:%s", userUUID)).RedisNil()

	isOptOut, err := repo.Is(userUUID, ctx)

	assert.NoError(t, err)
	assert.Equal(t, isOptOut, false)
}

func TestOptOutIsError(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := NewOptOutRepository(db)

	mock.Regexp().ExpectGet(fmt.Sprintf("opt-out:%s", userUUID)).SetErr(errors.New("error"))

	isOptOut, err := repo.Is(userUUID, ctx)

	assert.EqualError(t, err, "error")
	assert.Equal(t, isOptOut, false)
}
