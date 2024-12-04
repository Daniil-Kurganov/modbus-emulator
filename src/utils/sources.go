package utils

import "time"

var (
	ServerTCPHost     = "localhost"
	ServerTCPPort     = "1502"
	FinishDelayTime   = 3 * time.Second
	WorkMode          = "tcp"
	ModulePath        = `/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator`
	DumpDirectoryPath = `src/pcapng_files/main_files`
	Functions         = struct {
		CoilsRead          uint16
		DIRead             uint16
		HRRead             uint16
		IRRead             uint16
		CoilsSimpleWrite   uint16
		HRSimpleWrite      uint16
		CoilsMultipleWrite uint16
		HRMultipleWrite    uint16
	}{
		CoilsRead:          1,
		DIRead:             2,
		HRRead:             3,
		IRRead:             4,
		CoilsSimpleWrite:   5,
		HRSimpleWrite:      6,
		CoilsMultipleWrite: 15,
		HRMultipleWrite:    16,
	}
)
