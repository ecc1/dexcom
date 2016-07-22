package dexcom

// A Command specifies an operation to be performed by the Dexcom CGM receiver.
type Command byte

//go:generate stringer -type Command

const (
	NULL                          Command = 0
	ACK                           Command = 1
	NAK                           Command = 2
	INVALID_COMMAND               Command = 3
	INVALID_PARAM                 Command = 4
	INCOMPLETE_PACKET_RECEIVED    Command = 5
	RECEIVER_ERROR                Command = 6
	INVALID_MODE                  Command = 7
	PING                          Command = 10
	READ_FIRMWARE_HEADER          Command = 11
	READ_DATABASE_PARTITION_INFO  Command = 15
	READ_DATABASE_PAGE_RANGE      Command = 16
	READ_DATABASE_PAGES           Command = 17
	READ_DATABASE_PAGE_HEADER     Command = 18
	READ_TRANSMITTER_ID           Command = 25
	WRITE_TRANSMITTER_ID          Command = 26
	READ_LANGUAGE                 Command = 27
	WRITE_LANGUAGE                Command = 28
	READ_DISPLAY_TIME_OFFSET      Command = 29
	WRITE_DISPLAY_TIME_OFFSET     Command = 30
	READ_RTC                      Command = 31
	RESET_RECEIVER                Command = 32
	READ_BATTERY_LEVEL            Command = 33
	READ_SYSTEM_TIME              Command = 34
	READ_SYSTEM_TIME_OFFSET       Command = 35
	WRITE_SYSTEM_TIME             Command = 36
	READ_GLUCOSE_UNIT             Command = 37
	WRITE_GLUCOSE_UNIT            Command = 38
	READ_BLINDED_MODE             Command = 39
	WRITE_BLINDED_MODE            Command = 40
	READ_CLOCK_MODE               Command = 41
	WRITE_CLOCK_MODE              Command = 42
	READ_DEVICE_MODE              Command = 43
	ERASE_DATABASE                Command = 45
	SHUTDOWN_RECEIVER             Command = 46
	WRITE_PC_PARAMETERS           Command = 47
	READ_BATTERY_STATE            Command = 48
	READ_HARDWARE_BOARD_ID        Command = 49
	READ_FIRMWARE_SETTINGS        Command = 54
	READ_ENABLE_SETUP_WIZARD_FLAG Command = 55
	READ_SETUP_WIZARD_STATE       Command = 57
)
