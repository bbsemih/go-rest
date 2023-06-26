# gobank

gobank is a banking project written in Golang that allows you to create and manage bank accounts, record balance changes, and perform money transfers between accounts. It provides a set of RESTful HTTP APIs built using the Gin framework. The service uses PostgreSQL to store account information and transaction history. Docker is used for local development and GitHub Actions for running unit tests automatically.

## Features

* Serve both gRPC and HTTP requests by using grpc-gateway
* Create and manage bank accounts
* Record balance changes for each account
* Structured logs for both gRPC and HTTP by using zerolog
* Perform money transfers between two accounts
* RESTful HTTP APIs
* Store account information and transaction history in Postgres.
* Docker for local development
* GitHub Actions for automated unit tests

## TODO

* RabbitMQ
* Redis
* Chat Service


