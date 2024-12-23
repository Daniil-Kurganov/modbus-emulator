package src

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"net/http"
	"strconv"
	"strings"
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
		OneTimeEmulation bool   `json:"one_time_emulation"`
		StartTime        string `json:"start_time"`
		EndTime          string `json:"end_time"`
		CurrentTime      string `json:"current_time"`
	}
	serverData struct {
		ID       int             `json:"id"`
		Settings emulationServer `json:"settings"`
	}
	actualTime struct {
		ID         int    `json:"id"`
		ActualTime string `json:"actual_time"`
	}
	startEndTime struct {
		ID        int    `json:"id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
)

var (
	emulationServers struct {
		readWriteMutex sync.RWMutex
		serversData    []emulationServer
		rewindChannels []chan (int)
	}

	boolStringValues = map[string]bool{"true": true, "false": false}
	errorHeader      = "Error on HTTP-request"
)

func StartHTTPServer() {
	router := gin.Default()
	emulator := router.Group("/modbus-emulator")
	{
		settings := emulator.Group("settings")
		{
			settings.GET("", getSettings)
			settings.POST("emulation_mode", setEmulationMode)
		}
		time := emulator.Group("time")
		{
			time.GET("actual", getActualTime)
			time.GET("start&end", getStartEndTime)
			time.POST("rewind", rewindServers)
		}
		// emulator.GET("doc", func(gctx *gin.Context) {
		// 	gctx.Redirect(http.StatusPermanentRedirect,
		// 		fmt.Sprintf("http://%s:8080/modbus-emulator/docs/index.html", gctx.Request.Host),
		// 	)
		// })
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
	var leftBorder, rightBorder int
	serversData := getSettingsBuffer()
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		leftBorder, rightBorder = idInt, idInt+1
	} else {
		leftBorder, rightBorder = 0, len(serversData)
	}
	var response []serverData
	emulationServers.readWriteMutex.Lock()
	for currentID := range emulationServers.serversData[leftBorder:rightBorder] {
		emulationServers.serversData[currentID].OneTimeEmulation = flagValue
		response = append(response, serverData{
			ID:       currentID,
			Settings: emulationServers.serversData[currentID],
		})
	}
	emulationServers.readWriteMutex.Unlock()
	gctx.JSON(http.StatusOK, response)
}

func getActualTime(gctx *gin.Context) {
	serversData := getSettingsBuffer()
	var response []actualTime
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		emulationServers.readWriteMutex.RLock()
		response = append(response, actualTime{
			ID:         idInt,
			ActualTime: emulationServers.serversData[idInt].CurrentTime,
		})
		emulationServers.readWriteMutex.RUnlock()
	} else {
		emulationServers.readWriteMutex.RLock()
		for currentID, currentData := range emulationServers.serversData {
			response = append(response, actualTime{
				ID:         currentID,
				ActualTime: currentData.CurrentTime,
			})
		}
		emulationServers.readWriteMutex.RUnlock()
	}
	gctx.JSON(http.StatusOK, response)
}

func getStartEndTime(gctx *gin.Context) {
	serversData := getSettingsBuffer()
	var response []startEndTime
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		response = append(response, startEndTime{
			ID:        idInt,
			StartTime: serversData[idInt].StartTime,
			EndTime:   serversData[idInt].EndTime,
		})
	} else {
		for currentID, currentData := range serversData {
			response = append(response, startEndTime{
				ID:        currentID,
				StartTime: currentData.StartTime,
				EndTime:   currentData.EndTime,
			})
		}
	}
	gctx.JSON(http.StatusOK, response)
}

func rewindServers(gctx *gin.Context) {
	var timePointString string
	var ok bool
	if timePointString, ok = gctx.GetQuery("timepoint"); !ok {
		log.Printf("%s: missing required \"timepoint\" parameter", errorHeader)
		gctx.JSON(http.StatusBadRequest, gin.H{errorHeader: `missing required timepoint" parameter`})
		return
	}
	var err error
	var timepoint time.Time
	if timepoint, err = time.ParseInLocation(time.DateTime, timePointString, conf.DumpTimeLocation); err != nil {
		log.Printf("%s: invalid \"timepoint\" parameter", errorHeader)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: `invalid "timepoint" parameter`})
		return
	}
	serversData := getSettingsBuffer()
	var idString string
	if idString, ok = gctx.GetQuery("server_id"); ok {
		var err error
		var serverID int
		if serverID, err = strconv.Atoi(idString); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if serverID > len(serversData)-1 || serverID < 0 {
			log.Printf("Error on HTTP-request: \"server\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		var httpCode int
		if httpCode, ok, err = timeRewindTimepointCheck(serversData[serverID], timepoint); !ok {
			log.Printf("%s: %s", errorHeader, err)
			gctx.JSON(httpCode, gin.H{errorHeader: err.Error()})
			return
		}
		var transactionIndex int
		for currentIndex, currentTransaction := range History[serversData[serverID].DumpSocketsConfigData.RealSocket].Transactions {
			if timepoint.After(currentTransaction.TransactionTime) {
				transactionIndex = currentIndex
			}
		}
		emulationServers.readWriteMutex.RLock()
		emulationServers.rewindChannels[serverID] <- transactionIndex
		emulationServers.readWriteMutex.RUnlock()
		logString := fmt.Sprintf("Successfully rewinded %d server to %s timepoint", serverID, timepoint.String())
		log.Print(logString)
		gctx.JSON(http.StatusOK, gin.H{"Success": logString})
		return
	}
	// emulationServers.readWriteMutex.RLock()
	// for currentServerID, currentServerData := range emulationServers.serversData {
	// 	if serverID != -1 && serverID == currentServerID {

	// 	}
	// }
	// emulationServers.readWriteMutex.RUnlock()
	gctx.JSON(http.StatusOK, gin.H{"Success": timepoint.String()})
}

func getSettingsBuffer() []emulationServer {
	emulationServers.readWriteMutex.RLock()
	serversData := make([]emulationServer, len(emulationServers.serversData))
	copy(serversData, emulationServers.serversData)
	emulationServers.readWriteMutex.RUnlock()
	return serversData
}

func timeRewindTimepointCheck(serverData emulationServer, timepoint time.Time) (httpCode int, result bool, err error) {
	if !serverData.IsWorking {
		err = fmt.Errorf("current server isn't working")
		httpCode = http.StatusUnprocessableEntity
		return
	}
	if serverData.CurrentTime == "" {
		err = fmt.Errorf("current server isn't emulating data")
		httpCode = http.StatusUnprocessableEntity
		return
	}
	var startTime, endTime time.Time
	if startTime, err = time.ParseInLocation(time.DateTime, serverData.StartTime[:strings.Index(serverData.StartTime, " +")], conf.DumpTimeLocation); err != nil {
		err = fmt.Errorf("error on preprocessing start time: %s", err)
		httpCode = http.StatusInternalServerError
		return
	}
	if endTime, err = time.ParseInLocation(time.DateTime, serverData.EndTime[:strings.Index(serverData.EndTime, " +")], conf.DumpTimeLocation); err != nil {
		err = fmt.Errorf("error on preprocessing end time: %s", err)
		httpCode = http.StatusInternalServerError
		return
	}
	if !(timepoint.After(startTime) && timepoint.Before(endTime)) {
		err = fmt.Errorf("\"timepoint\" must be between %s and %s", startTime.String(), endTime.String())
		httpCode = http.StatusInternalServerError
		return
	}
	result = true
	return
}
