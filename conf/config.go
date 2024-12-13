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
			HostAddress: "192.168.1.29",
			PortAddress: "502",
		},
		1503: {
			HostAddress: "192.168.1.31",
			PortAddress: "502",
		},
	}
	ServerTCPHost     = "127.0.0.1"
	FinishDelayTime   = 3 * time.Second
	WorkMode          = "rtu_over_tcp"
	ModulePath        = `/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator`
	DumpDirectoryPath = `pcapng_files/main_files`
)
