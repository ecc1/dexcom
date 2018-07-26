package dexcom

var (
	// Ensure that *bleConn implements the Connection interface.
	_ Connection = (*bleConn)(nil)

	// Ensure that *usbConn implements the Connection interface.
	_ Connection = (*usbConn)(nil)
)
