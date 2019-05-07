#!/bin/bash

go run github.com/golang/mock/mockgen -package=mock -source=$1.go -destination=../../pkg/mock/$1.go
