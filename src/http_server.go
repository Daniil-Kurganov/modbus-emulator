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
	emulationServer struct {
		IsWorking bool `json:"is_working"`
		conf.DumpSocketsConfigData
		OneTimeEmulation bool      `json:"one_time_emulation"`
		StartTime        time.Time `json:"start_time"`
		EndTime          time.Time `json:"end_time"`
		CurrentTime      time.Time `json:"current_time"`
	}
	serverData struct {
		ID       int             `json:"id"`
		Settings emulationServer `json:"settings"`
	}
)

var (
	emulationServers struct {
		readWriteMutex sync.RWMutex
		serversData    []emulationServer
	}

	boolStringValues = map[string]bool{"true": true, "false": false}
	errorHeader      = "Error on HTTP-request"
)

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
			settings.GET("", getSettings)
			settings.POST("emulation_mode", setEmulationMode)
		}
		emulator.GET("doc", func(gctx *gin.Context) {
			gctx.Redirect(http.StatusPermanentRedirect,
				fmt.Sprintf("http://%s:8080/modbus-emulator/docs/index.html", gctx.Request.Host),
			)
		})
		emulator.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.Run(conf.ServerHTTPServesocket)
}

func getSettings(gctx *gin.Context) {
	var response []serverData
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		serversData := getSettingsBuffer()
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		response = append(response, serverData{
			ID:       idInt,
			Settings: serversData[idInt],
		})
	} else {
		for currentID, currentSetting := range getSettingsBuffer() {
			response = append(response,
				serverData{
					ID:       currentID,
					Settings: currentSetting,
				},
			)
		}
	}
	gctx.JSON(http.StatusOK, response)
}

func setEmulationMode(gctx *gin.Context) {
	var mode string
	var ok bool
	if mode, ok = gctx.GetQuery("one-time"); !ok {
		errorLog := "missig \"one-time\" parameter"
		log.Printf("%s: %s", errorHeader, errorLog)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: errorLog})
		return
	}
	var flagValue bool
	if flagValue, ok = boolStringValues[mode]; !ok {
		errorLog := "invalid \"one-time\" parameter (must be \"true\" or \"false\")"
		log.Printf("%s: %s", errorHeader, errorLog)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: errorLog})
		return
	}
	var response []serverData
	emulationServers.readWriteMutex.Lock()
	for currentID := range emulationServers.serversData {
		emulationServers.serversData[currentID].OneTimeEmulation = flagValue
		response = append(response, serverData{
			ID:       currentID,
			Settings: emulationServers.serversData[currentID],
		})
	}
	emulationServers.readWriteMutex.Unlock()
	gctx.JSON(http.StatusOK, response)
}

// func getSettings(gctx *gin.Context) {
// 	var err error
// 	var id int
// 	if id, err = strconv.Atoi(gctx.Param("server")); err != nil {
// 		log.Printf("Error on HTTP-request: invalid \"server\" parameter: %s", err)
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{"Invalid \"server\" parameter": err.Error()})
// 		return
// 	}
// 	serversData := getSettingsBuffer()
// 	if id > len(serversData)-1 || id < 0 {
// 		log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
// 		return
// 	}
// 	response := serverData{
// 		ID:       id,
// 		Settings: serversData[id],
// 	}
// 	gctx.JSON(http.StatusOK, response)
// }

// func setEmulationMode(gctx *gin.Context) {
// 	var err error
// 	var id int
// 	if id, err = strconv.Atoi(gctx.Param("server")); err != nil {
// 		log.Printf("%s: invalid \"server\" parameter: %s", errorHeader, err)
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{"Invalid \"server\" parameter": err.Error()})
// 		return
// 	}
// 	serversData := getSettingsBuffer()
// 	if id > len(serversData)-1 || id < 0 {
// 		log.Printf("%s: \"server\" parameter must be in range [0:%d]", errorHeader, len(serversData))
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
// 		return
// 	}
// 	var mode string
// 	var ok bool
// 	if mode, ok = gctx.GetQuery("one-time"); !ok {
// 		errorLog := "missig \"one-time\" parameter"
// 		log.Printf("%s: %s", errorHeader, errorLog)
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: errorLog})
// 		return
// 	}
// 	var flagValue bool
// 	if flagValue, ok = boolStringValues[mode]; !ok {
// 		errorLog := "invalid \"one-time\" parameter (must be \"true\" or \"false\")"
// 		log.Printf("%s: %s", errorHeader, errorLog)
// 		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: errorLog})
// 		return
// 	}
// 	emulationServers.readWriteMutex.Lock()
// 	emulationServers.serversData[id].OneTimeEmulation = flagValue
// 	response := serverData{
// 		ID:       id,
// 		Settings: emulationServers.serversData[id],
// 	}
// 	gctx.JSON(http.StatusOK, response)
// }

func getSettingsBuffer() []emulationServer {
	emulationServers.readWriteMutex.RLock()
	serversData := make([]emulationServer, len(emulationServers.serversData))
	copy(serversData, emulationServers.serversData)
	emulationServers.readWriteMutex.RUnlock()
	return serversData
}
