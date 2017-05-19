export GOPATH=$(shell pwd)

all:
	go build gimbal_control

deps:
	go get github.com/jacobsa/go-serial/serial
	go get github.com/adammck/dynamixel
