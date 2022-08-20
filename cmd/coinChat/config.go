package main

import "coinChat/pkg/repository"

type Config struct {
	Service
	repository.ConfigPostgres
}

type Service struct {
	Env  string
	Name string
}
