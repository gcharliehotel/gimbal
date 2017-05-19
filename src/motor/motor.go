package motor

import (
	"encoding/binary"
	"fmt"
	"github.com/adammck/dynamixel/iface"
	"os"
	"time"
)

type Motor struct {
	proto iface.Protocol
	Id    int
}

const regAddrMoving = 122              // byte
const regAddrGoalPosition = 116        // 4 bytes
const regAddrTorqueEnable = 64         // byte
const regAddrProfileAcceleration = 108 // 4 bytes
const regAddrProfileVelocity = 112     // 4 bytes

func New(proto iface.Protocol, id int) *Motor {
	return &Motor{proto: proto, Id: id}
}

func (motor *Motor) Ping() error {
	return motor.proto.Ping(motor.Id)
}

func (motor *Motor) SetGoalPosition(position int) {
	fmt.Printf("setting goal position of %d to %d\n", motor.Id, position)
	motor.writeUint32(regAddrGoalPosition, uint32(position))
}

func (motor *Motor) SetGoalPositionAndWaitForStop(position int) {
	motor.SetGoalPosition(position)
	motor.WaitForStop()
}

func (motor *Motor) SetTorqueEnable(enable bool) {
	if enable {
		motor.writeByte(regAddrTorqueEnable, 1)
	} else {
		motor.writeByte(regAddrTorqueEnable, 0)
	}
}

func (motor *Motor) SetProfileAcceleration(acceleration int) {
	motor.writeUint32(regAddrProfileAcceleration, uint32(acceleration))
}

func (motor *Motor) SetProfileVelocity(velocity int) {
	motor.writeUint32(regAddrProfileVelocity, uint32(velocity))
}

func (motor *Motor) WaitForStop() {
	fmt.Printf("waiting for %d to stop\n", motor.Id)
	for {
		moving := motor.readByte(regAddrMoving)
		if moving == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (motor *Motor) writeData(addr int, bytes []byte) error {
	return motor.proto.WriteData(motor.Id, addr, bytes, true)
}

func (motor *Motor) readData(addr int, len int) ([]byte, error) {
	return motor.proto.ReadData(motor.Id, addr, len)
}

func (motor *Motor) writeByte(addr int, value byte) {
	bytes := []byte{value}
	err := motor.writeData(addr, bytes)
	if err != nil {
		fmt.Printf("failed to write id=%d, addr=%d, value=%d: %s", motor.Id, addr, value, err)
		os.Exit(1)
	}
}

func (motor *Motor) writeUint32(addr int, value uint32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(value))
	err := motor.writeData(addr, bytes)
	if err != nil {
		fmt.Printf("failed to write id=%d, addr=%d, value=%d: %s", motor.Id, addr, value, err)
		os.Exit(1)
	}
}

func (motor *Motor) readByte(addr int) byte {
	len := 1
	data, err := motor.readData(addr, len)
	if err != nil {
		fmt.Printf("failed to read id=%d, addr=%d, len=%d: %s", motor.Id, addr, len, err)
		os.Exit(1)
	}
	return data[0]
}
