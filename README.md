# Go software for CAN-Open and RR IO card testing

The software uses the PEAK drivers (https://www.peak-system.com/quick/DrvSetup) and they
must be installed on the comuter running the test programs.

The software supports control of power supplies over USB/RS232. It has only been tested 
for TTi CPX4000DP, but uses standard commands to set voltage/current, and it should be
easy to modify it for different power supplies. If no power supply is found, 
manual instructions to the user is used instead to turn power on/off. 

## Installation and setup

The following dependencies are used, and will be loaded by the following commands:
(The last one is used only for ```go test```)
```
	go get github.com/gookit/color
	go get github.com/stretchr/testify
```

To run the program directly, use
```
    go run ./rraio16/
```

To build the exe file, just use type
```
    go build ./rraio16/
```

This will make a file rraio16.exe in the same directory. This is a standalone exe file that
can be copied anywhere, and have no dependencies except for the peak can-bus drivers.

The software can be compiled on a Linux machine, but the Peak can drivers have not 
been tested on Linux. A new peak.go module for Linux will be needed.

## Flags

```
  -help
        Will print this test info  
  -node int
        Node number to test. Uses 11 as default. (default 11)
  -power-port string
        Name of com-port used to control the power supply, defaults to device with highest com-port number
  -subtest int
        Set to the subtest number that should be executed. Set to zero to run all tests
  -test-power
        Check if a power supply is connected and verify connection
  -toggle-power
        Set to false to disable power off-on at start of test. This also skips test of bootup messages. (default true)
```

## Testing RRAIO16 cards
### Setup
The test program assumes that a RRAIO16 card is connected in the following way.
* A TTiCPX4000 power supply connected via USB will be used if present. Fallback is manual settings.
* The can bus is connected via a Peak USB adapter as device 1
* Power supply should be 20.00V, and must be connected to terminal 1-2 and 9-10
* Channel 9 and 13 is connected together
* Channel 10 and 14 is connected together
* Channel 11 and 15 is connected together
* Channel 12 and 16 is connected together
* Channel 1 and 3 is connected together
* Channel 2 and 4 is connected together
* Channel 1 and 7 is connected  together via 1 kohm resistor
* Channel 2 and 8 is connected  together via 1 kohm resistor

### List of subtests implemented
1. Verify mandatory objects from 0x1000 to 0x1018
2. Verify tx pdo parameters at 0x1800
3. Verify tx pdo mapping at 0x1A00
4. Verify pdo parameters at 0x1400
5. Verify pdo mapping at 0x1600
6. Verify heartbeat operation"
7. Verify emergency message when starting pdos
8. Testing rx pdo operation - transfer output values to card
9. Verify tx pdos
10. Readback digital outputs
11. Testing analog inputs
12. Frequency measurement

