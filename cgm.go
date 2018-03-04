/*
Package dexcom provides functions to access a Dexcom G4 Share
CGM receiver over a BLE or USB connection.

Based on the Python version at github.com/bewest/decoding-dexcom
*/
package dexcom

// Connection is the interface satisfied by a CGM connection.
type Connection interface {
	Send([]byte) error
	Receive([]byte) error
	Close()
}

// CGM represents a CGM connection.
type CGM struct {
	Connection
	err error
}

// Open first attempts to open a USB connection;
// if that fails it tries a BLE connection.
func Open() *CGM {
	conn, err := OpenUSB()
	if err == nil {
		return &CGM{Connection: conn}
	}
	conn, err = OpenBLE()
	return &CGM{Connection: conn, err: err}
}

// Error returns the error state of the CGM.
func (cgm *CGM) Error() error {
	return cgm.err
}

// SetError sets the error state of the CGM.
func (cgm *CGM) SetError(err error) {
	cgm.err = err
}
