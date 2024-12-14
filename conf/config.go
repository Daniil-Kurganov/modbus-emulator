package conf

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type ServerSocketData struct {
	HostAddress string
	PortAddress string
	WorkMode    string
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

	Ports             map[string]ServerSocketData
	ServerTCPHost     = "0.0.0.0"
	FinishDelayTime   = 3 * time.Second
	ModulePath        = `/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator`
	DumpDirectoryPath = `pcapng_files/main_files`
	DumpFileName      = "main"
)

func init() {
	log.SetFlags(0)
	Ports = make(map[string]ServerSocketData)
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", ModulePath), &Ports); err != nil {
		log.Fatalf("Error on unmarshaling configuration: %s", err)
	}
}
