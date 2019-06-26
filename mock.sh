#!/bin/bash

go run github.com/golang/mock/mockgen -package=mock -source=$GOFILE -destination=../../pkg/mock/$GOFILE
