package conf

import "time"

type ServerSocket struct {
	HostAddress string
	PortAddress string
}

var (
	Functions = struct {
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

	Ports = map[uint16]ServerSocket{
		1502: {
			HostAddress: "127.0.0.1",
			PortAddress: "1502",
		},
		1503: {
			HostAddress: "127.0.0.1",
			PortAddress: "1503",
		},
	}
	ServerTCPHost     = "localhost"
	FinishDelayTime   = 3 * time.Second
	WorkMode          = "tcp"
	ModulePath        = `/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator`
	DumpDirectoryPath = `pcapng_files/main_files/multiple_ports`
)
