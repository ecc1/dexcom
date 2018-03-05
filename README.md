# dexcom

The `dexcom` package provides functions to connect to Dexcom G4 and
Dexcom G4 Share continuous glucose monitor (CGM) systems.
It connects by either USB or Bluetooth LE to the Dexcom receiver.

### Wireless connection (Bluetooth LE; G4 Share only)

For BLE connections, the serial number of the Dexcom receiver
must be passed in the `DEXCOM_CGM_ID` environment variable.

For the initial connection to the receiver, use the `Settings > Share` menu
to `Forget Device`, then turn Share back on.
The `cmd/g4ping` program can be used to confirm connection with the receiver.

### Wired connection (USB serial)

For serial connections, the host must have the `cdc_acm` kernel module loaded.
Look for the creation of `/dev/ttyACM0` in the system log
when the receiver is attached.

### Utility programs

The `cmd` directory contains some simple utility programs:

* `g4ping` pings the receiver (first connecting if necessary)
  and exits with a success or failure status.
* `glucose` retrieves CGM data and prints it in various formats.
* `backfill` finds gaps in
 [Nightscout](https://github.com/nightscout/cgm-remote-monitor) CGM data,
 retrieves the missing data from the receiver, and uploads it.
 Note that a USB connection works much faster for gaps
 that are hours or days in the past, and can be done from any Linux machine,
 not just an [OpenAPS](https://github.com/openapsopenaps) rig.
* `g4setclock` sets the receiver's date and time.
* `g4update` retrieves CGM data, with options to update a local JSON file
 and upload to [Nightscout.](https://github.com/nightscout/cgm-remote-monitor)

### Documentation

<https://godoc.org/github.com/ecc1/dexcom>
