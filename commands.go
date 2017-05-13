package dexcom

// Command represents a Dexcom CGM receiver command.
type Command byte

//go:generate stringer -type Command

// Dexcom G4 receiver commands.
const (
	Null                       Command = 0
	Ack                        Command = 1
	Nak                        Command = 2
	InvalidCommand             Command = 3
	InvalidParam               Command = 4
	IncompletePacketReceived   Command = 5
	ReceiverError              Command = 6
	InvalidMode                Command = 7
	Ping                       Command = 10
	ReadFirmwareHeader         Command = 11
	ReadDatabasePartitionInfo  Command = 15
	ReadDatabasePageRange      Command = 16
	ReadDatabasePages          Command = 17
	ReadDatabasePageHeader     Command = 18
	ReadTransmitterID          Command = 25
	WriteTransmitterID         Command = 26
	ReadLanguage               Command = 27
	WriteLanguage              Command = 28
	ReadDisplayTimeOffset      Command = 29
	WriteDisplayTimeOffset     Command = 30
	ReadRTC                    Command = 31
	ResetReceiver              Command = 32
	ReadBatteryLevel           Command = 33
	ReadSystemTime             Command = 34
	ReadSystemTimeOffset       Command = 35
	WriteSystemTime            Command = 36
	ReadGlucoseUnits           Command = 37
	WriteGlucoseUnits          Command = 38
	ReadBlindMode              Command = 39
	WriteBlindMode             Command = 40
	ReadClockMode              Command = 41
	WriteClockMode             Command = 42
	ReadDeviceMode             Command = 43
	EraseDatabase              Command = 45
	ShutdownReceiver           Command = 46
	WriteSoftwareParameters    Command = 47
	ReadBatteryState           Command = 48
	ReadHardwareID             Command = 49
	ReadFirmwareSettings       Command = 54
	ReadEnableSetupWizardFlag  Command = 55
	ReadSetupWizardState       Command = 57
	ReadChargerCurrentSetting  Command = 59
	WriteChargerCurrentSetting Command = 60
)
