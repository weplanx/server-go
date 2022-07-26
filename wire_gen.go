// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/api/dsl"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/sessions"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
)

// Injectors from wire.go:

func OkLetsGo(value *common.Values) (*server.Hertz, error) {
	client, err := bootstrap.UseMongoDB(value)
	if err != nil {
		return nil, err
	}
	database := bootstrap.UseDatabase(client, value)
	redisClient, err := bootstrap.UseRedis(value)
	if err != nil {
		return nil, err
	}
	conn, err := bootstrap.UseNats(value)
	if err != nil {
		return nil, err
	}
	service := &values.Service{
		Values: value,
		Db:     database,
		Redis:  redisClient,
		Nats:   conn,
	}
	usersService := &users.Service{
		Db:    database,
		Redis: redisClient,
	}
	indexService := &index.Service{
		Values: value,
		Mongo:  client,
		Db:     database,
		Redis:  redisClient,
		Users:  usersService,
	}
	apiAPI := &api.API{
		Values:        value,
		ValuesService: service,
		IndexService:  indexService,
		UsersService:  usersService,
	}
	controller := &index.Controller{
		AppService: indexService,
	}
	valuesController := &values.Controller{
		ValuesService: service,
	}
	sessionsService := &sessions.Service{
		Values: value,
		Redis:  redisClient,
	}
	sessionsController := &sessions.Controller{
		SessionsService: sessionsService,
	}
	dslService := &dsl.Service{
		Db: database,
	}
	dslController := &dsl.Controller{
		DslService: dslService,
	}
	hertz, err := api.Routes(apiAPI, controller, valuesController, sessionsController, dslController)
	if err != nil {
		return nil, err
	}
	return hertz, nil
}
