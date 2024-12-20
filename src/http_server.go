package src

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "modbus-emulator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	workServer struct {
		conf.DumpSocketsConfigData
		OneTimeEmulation bool      `json:"one_time_emulation"`
		StartTime        time.Time `json:"start_time"`
		EndTime          time.Time `json:"end_time"`
		CurrentTime      time.Time `json:"current_time"`
	}
	serverData struct {
		ID       int        `json:"id"`
		Settings workServer `json:"settings"`
	}
)

var workServers struct {
	readWriteMutex sync.RWMutex
	serversData    []workServer
}

func StartHTTPServer() {
	router := gin.Default()
	emulator := router.Group("/modbus-emulator")
	{
		time := emulator.Group("time")
		{
			time.GET("actual")            // info about current time of all working emulations servers
			time.GET("start&end")         // info about starting and ending times all working emulations servers
			time.POST("rewind")           // rewind all working servers on a specified time
			time.GET(":server/actual")    // info of current time of server
			time.GET(":server/start&end") // info about starting and ending of server
			time.POST(":server/rewind")   // rewind server on a specified time
		}
		settings := emulator.Group("settings")
		{
			settings.GET("", getAllSettings)
			settings.POST("emulation_mode") // set emulation mode for all server
			settings.GET(":server", getSettings)
			settings.POST(":server/emulation_mode") // set emulation mode for server
		}
		emulator.GET("doc", func(gctx *gin.Context) {
			gctx.Redirect(http.StatusPermanentRedirect,
				fmt.Sprintf("http://%s/modbus-emulator/docs/index.html", conf.ServerHTTPServesocket),
			)
		})
		emulator.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.Run(conf.ServerHTTPServesocket)
}

func getAllSettings(gctx *gin.Context) {
	var response []serverData
	for currentID, currentSetting := range getSettingsBuffer() {
		response = append(response,
			serverData{
				ID:       currentID,
				Settings: currentSetting,
			},
		)
	}
	gctx.JSON(http.StatusOK, response)
}

func getSettings(gctx *gin.Context) {
	var err error
	var id int
	if id, err = strconv.Atoi(gctx.Param("server")); err != nil {
		log.Printf("Error on HTTP-request: invalid \"server\" parameter: %s", err)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{"Invalid \"server\" parameter": err.Error()})
		return
	}
	serversData := getSettingsBuffer()
	if id > len(serversData)-1 || id < 0 {
		log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
		return
	}
	response := serverData{
		ID:       id,
		Settings: serversData[id],
	}
	gctx.JSON(http.StatusOK, response)
}

func getSettingsBuffer() []workServer {
	workServers.readWriteMutex.RLock()
	serversData := make([]workServer, len(workServers.serversData))
	copy(serversData, workServers.serversData)
	workServers.readWriteMutex.RUnlock()
	return serversData
}
