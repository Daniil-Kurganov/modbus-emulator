package src

import (
	"fmt"
	"modbus-emulator/conf"
	"net/http"
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
			settings.POST("emulation_mode")         // set emulation mode for all server
			settings.GET(":server")                 // settings of server
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
	workServers.readWriteMutex.RLock()
	serversData := make([]workServer, len(workServers.serversData))
	copy(serversData, workServers.serversData)
	workServers.readWriteMutex.RUnlock()
	for currentID, currentSetting := range serversData {
		response = append(response,
			serverData{
				ID:       currentID,
				Settings: currentSetting,
			},
		)
	}
	gctx.JSON(http.StatusOK, response)
}
