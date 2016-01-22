package dexcom

type Connection interface {
	Send([]byte) error
	Receive([]byte) error
}

var conn Connection

func Open() error {
	var err error
	conn, err = OpenUSB()
	return err
}
