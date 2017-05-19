package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adammck/dynamixel/network"
	proto2 "github.com/adammck/dynamixel/protocol/v2"
	"github.com/jacobsa/go-serial/serial"

	"motor"
)

var (
	portName     = flag.String("port", "/dev/ttyUSB0", "the serial port path")
	acceleration = flag.Int("acceleration", 50, "the profile acceleration to set")
	velocity     = flag.Int("velocity", 50, "the profile velocity to set")
	debug        = flag.Bool("debug", false, "show serial traffic")
	unload       = flag.Bool("unload", false, "move to unload device")
	home         = flag.Bool("home", true, "home the device")
	run          = flag.Bool("run", false, "run the profile")
)

func main() {
	flag.Parse()

	options := serial.OpenOptions{
		PortName:              *portName,
		BaudRate:              57600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	serial, err := serial.Open(options)
	if err != nil {
		log.Fatalf("error opening serial port: %v\n", err)
	}

	network := network.New(serial)
	if *debug {
		network.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	network.Timeout = 50 * time.Millisecond

	network.Flush()

	proto := proto2.New(network)

	yaw := motor.New(proto, 1)
	pitch := motor.New(proto, 2)
	roll := motor.New(proto, 3)

	all_motors := []*motor.Motor{yaw, pitch, roll}

	for _, motor := range all_motors {
		err = motor.Ping()
		if err != nil {
			fmt.Printf("ping error: %s\n", err)
			os.Exit(1)
		}

		motor.SetTorqueEnable(true)
		motor.SetProfileVelocity(*velocity)
		motor.SetProfileAcceleration(*acceleration)
	}

	center := 2048
	deg45 := 512
	deg90 := 1024
	deg180 := 2048

	if *home {
		for _, motor := range all_motors {
			motor.SetGoalPositionAndWaitForStop(center)
		}
	}

	if *run {
		yaw.SetGoalPositionAndWaitForStop(center - deg90)
		yaw.SetGoalPositionAndWaitForStop(center + deg90)
		yaw.SetGoalPositionAndWaitForStop(center)

		pitch.SetGoalPositionAndWaitForStop(center - deg45)
		pitch.SetGoalPositionAndWaitForStop(center + deg45)
		pitch.SetGoalPositionAndWaitForStop(center)

		roll.SetGoalPositionAndWaitForStop(center - deg90)
		roll.SetGoalPositionAndWaitForStop(center + deg90)
		roll.SetGoalPositionAndWaitForStop(center)
	}

	if *unload {
		pitch.SetGoalPositionAndWaitForStop(center - deg180)
	}
}
