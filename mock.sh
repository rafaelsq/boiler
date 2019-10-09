#!/bin/bash

mockgen -package=mock -source=$GOFILE -destination=../../pkg/mock/$GOFILE
