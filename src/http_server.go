package src

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modbus-emulator/docs"

	ms "github.com/Daniil-Kurganov/modbus-server"
	"github.com/gin-gonic/gin"
	reuse "github.com/libp2p/go-reuseport"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/maps"
)

type (
	emulationServerSettings struct {
		IsWorking   bool `json:"is_working"`
		IsEmulating bool `json:"is_emulating"`
		conf.DumpSocketsConfigData
		OneTimeEmulation bool   `json:"one_time_emulation"`
		StartTime        string `json:"start_time"`
		EndTime          string `json:"end_time"`
		CurrentTime      string `json:"current_time"`
	}
	settingsResponse struct {
		ID       int                     `json:"id"`
		Settings emulationServerSettings `json:"settings"`
	}
	actualTimeResponse struct {
		ID         int    `json:"id"`
		ActualTime string `json:"actual_time"`
	}
	startEndTimeRespoonse struct {
		ID        int    `json:"id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	rewindResponse struct {
		ID              int    `json:"id"`
		Error           string `json:"error"`
		SettedTimepoint string `json:"setted_timepoint"`
	}
	slaveResponse struct {
		ServerID        int    `json:"server_id"`
		AnswerwedSlaves []int8 `json:"answered_slaves"`
	}
	emulationControlResponse struct {
		ID          int    `json:"id"`
		Error       string `json:"error"`
		IsEmulating bool   `json:"is_emulating"`
	}
)

var (
	emulationServers struct {
		readWriteMutex           sync.RWMutex
		serversData              []emulationServerSettings
		servers                  []*ms.Server
		rewindChannels           []chan (int)
		emulationControlChannels []chan (bool)
	}

	boolStringValues = map[string]bool{"true": true, "false": false, "start": true, "stop": false}
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
			settings.POST("slave_answer", setSlaveState)
			settings.POST("emulation_control", controlEmulation)
		}
		time := emulator.Group("time")
		{
			time.GET("actual", getActualTime)
			time.GET("start&end", getStartEndTime)
			time.POST("rewind_emulation", rewindServersEmulation)
		}
		emulator.GET("/", func(gctx *gin.Context) {
			gctx.Redirect(http.StatusPermanentRedirect,
				fmt.Sprintf("http://%s/modbus-emulator/docs/index.html", gctx.Request.Host),
			)
		})
		emulator.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	var listener net.Listener
	var err error
	if listener, err = reuse.Listen("tcp", conf.ServerHTTPServesocket); err != nil {
		log.Fatalf("Error on creating listener: %s", err)
	}
	if err = router.RunListener(listener); err != nil {
		log.Fatalf("Error on starting HTTP-server: %s", err)
	}
}

func getSettings(gctx *gin.Context) {
	var response []settingsResponse
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
		response = append(response, settingsResponse{
			ID:       idInt,
			Settings: serversData[idInt],
		})
	} else {
		for currentID, currentSetting := range getSettingsBuffer() {
			response = append(response,
				settingsResponse{
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
	var response []settingsResponse
	emulationServers.readWriteMutex.Lock()
	for currentID := range emulationServers.serversData[leftBorder:rightBorder] {
		emulationServers.serversData[currentID].OneTimeEmulation = flagValue
		response = append(response, settingsResponse{
			ID:       currentID,
			Settings: emulationServers.serversData[currentID],
		})
	}
	emulationServers.readWriteMutex.Unlock()
	gctx.JSON(http.StatusOK, response)
}

func setSlaveState(gctx *gin.Context) {
	var err error
	serversData := getSettingsBuffer()
	var serverID, slaveID int
	var workMode bool
	var id, mode string
	var ok bool
	if id, ok = gctx.GetQuery("server_id"); !ok {
		err = fmt.Errorf("missed required \"server_id\" parameter")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusBadRequest, gin.H{errorHeader: err.Error()})
		return
	}
	if serverID, err = strconv.Atoi(id); err != nil {
		log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
		return
	}
	if serverID > len(serversData)-1 || serverID < 0 {
		log.Printf("Error on HTTP-request: \"server_id\" parameter must be in range [0:%d]", len(serversData))
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
		return
	}
	if id, ok = gctx.GetQuery("slave_id"); !ok {
		err = fmt.Errorf("missed required \"slave_id\" parameter")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusBadRequest, gin.H{errorHeader: err.Error()})
		return
	}
	if slaveID, err = strconv.Atoi(id); err != nil {
		err = fmt.Errorf("invalid \"slave_id\" parameter: %s", err)
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
		return
	}
	if mode, ok = gctx.GetQuery("answer_mode"); !ok {
		err = fmt.Errorf("missed required \"answer_mode\" parameter")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusBadRequest, gin.H{errorHeader: err.Error()})
		return
	}
	if workMode, ok = boolStringValues[mode]; !ok || !slices.Contains([]string{"start", "stop"}, mode) {
		err = fmt.Errorf("invalid \"answer_mode\" parameter (must be \"start\" or \"stop\")")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
		return
	}
	emulationServers.readWriteMutex.RLock()
	defer emulationServers.readWriteMutex.RUnlock()
	if workMode {
		if err = emulationServers.servers[serverID].SlaveStartResponse(uint8(slaveID)); err != nil {
			log.Printf("%s: %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
			return
		}
	}
	if !workMode {
		if err = emulationServers.servers[serverID].SlaveStopResponse(uint8(slaveID)); err != nil {
			log.Printf("%s: %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
			return
		}
	}
	slavesAnswered := maps.Keys(emulationServers.servers[serverID].Slaves)
	slavesStoped := emulationServers.servers[serverID].SlavesStoppedResponse
	var slavesResponse []int8
	for _, currentSlave := range slavesAnswered {
		if !slices.Contains(slavesStoped, currentSlave) {
			slavesResponse = append(slavesResponse, int8(currentSlave))
		}
	}
	response := slaveResponse{ServerID: serverID}
	response.AnswerwedSlaves = slavesResponse
	gctx.JSON(http.StatusOK, response)
}

func controlEmulation(gctx *gin.Context) {
	var err error
	var flag string
	var ok bool
	if flag, ok = gctx.GetQuery("emulation_switch"); !ok {
		err = fmt.Errorf("missed required \"emulation_switch\" parameter")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusBadRequest, gin.H{errorHeader: err.Error()})
		return
	}
	var isEmulating bool
	if isEmulating, ok = boolStringValues[flag]; !ok || !slices.Contains([]string{"start", "stop"}, flag) {
		err = fmt.Errorf("invalid \"emulation_switch\" parameter (must be \"start\" or \"stop\")")
		log.Printf("%s: %s", errorHeader, err)
		gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
		return
	}
	var id string
	serversData := getSettingsBuffer()
	var response []emulationControlResponse
	if id, ok = gctx.GetQuery("server_id"); ok {
		var serverID int
		if serverID, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if serverID > len(serversData)-1 || serverID < 0 {
			log.Printf("Error on HTTP-request: \"server_id\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		if serversData[serverID].IsEmulating == isEmulating {
			err = fmt.Errorf("invalid \"emulation_switch\" parameter: server is already in requested state")
			log.Printf("%s: %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
			return
		}
		emulationServers.readWriteMutex.RLock()
		defer emulationServers.readWriteMutex.RUnlock()
		select {
		case emulationServers.emulationControlChannels[serverID] <- true:
			response = append(response, emulationControlResponse{
				ID:          serverID,
				IsEmulating: isEmulating,
			})
			gctx.JSON(http.StatusOK, response)
			return
		case <-time.After(time.Second):
			err = fmt.Errorf("invalid \"emulation_switch\" parameter: server couldn't process state (emulation isn't initialized)")
			log.Printf("%s: %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{errorHeader: err.Error()})
			return
		}
	}
	for currentID, currentData := range serversData {
		currentResponse := emulationControlResponse{ID: currentID}
		if currentData.IsEmulating == isEmulating {
			err = fmt.Errorf("invalid \"emulation_switch\" parameter: server is already in requested state")
			log.Printf("%s: %s", errorHeader, err)
			currentResponse.Error = err.Error()
			currentResponse.IsEmulating = currentData.IsEmulating
		} else {
			emulationServers.readWriteMutex.RLock()
			defer emulationServers.readWriteMutex.RUnlock()
			select {
			case emulationServers.emulationControlChannels[currentID] <- true:
				currentResponse.IsEmulating = isEmulating
			case <-time.After(time.Second):
				err = fmt.Errorf("invalid \"emulation_switch\" parameter: server couldn't process state (emulation isn't initialized)")
				log.Printf("%s: %s", errorHeader, err)
				currentResponse.Error = err.Error()
				currentResponse.IsEmulating = currentData.IsEmulating
			}
		}
		response = append(response, currentResponse)
	}
	gctx.JSON(http.StatusOK, response)
}

func getActualTime(gctx *gin.Context) {
	serversData := getSettingsBuffer()
	var response []actualTimeResponse
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server_id\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		emulationServers.readWriteMutex.RLock()
		response = append(response, actualTimeResponse{
			ID:         idInt,
			ActualTime: emulationServers.serversData[idInt].CurrentTime,
		})
		emulationServers.readWriteMutex.RUnlock()
	} else {
		emulationServers.readWriteMutex.RLock()
		for currentID, currentData := range emulationServers.serversData {
			response = append(response, actualTimeResponse{
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
	var response []startEndTimeRespoonse
	if id, ok := gctx.GetQuery("server_id"); ok {
		var idInt int
		var err error
		if idInt, err = strconv.Atoi(id); err != nil {
			log.Printf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`Invalid "server_id" parameter`: err.Error()})
			return
		}
		if idInt > len(serversData)-1 || idInt < 0 {
			log.Printf("Error on HTTP-request: \"server_id\" parameter must be in range [0:%d]", len(serversData))
			gctx.JSON(http.StatusUnprocessableEntity, gin.H{`"server" parameter must be in range`: fmt.Sprintf("[0:%d]", len(serversData))})
			return
		}
		response = append(response, startEndTimeRespoonse{
			ID:        idInt,
			StartTime: serversData[idInt].StartTime,
			EndTime:   serversData[idInt].EndTime,
		})
	} else {
		for currentID, currentData := range serversData {
			response = append(response, startEndTimeRespoonse{
				ID:        currentID,
				StartTime: currentData.StartTime,
				EndTime:   currentData.EndTime,
			})
		}
	}
	gctx.JSON(http.StatusOK, response)
}

func rewindServersEmulation(gctx *gin.Context) {
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
	var response []rewindResponse
	if idString, ok = gctx.GetQuery("server_id"); ok {
		var err error
		var serverID int
		responseValue := rewindResponse{}
		if serverID, err = strconv.Atoi(idString); err != nil {
			err = fmt.Errorf("%s: invalid \"server_id\" parameter - %s", errorHeader, err)
			log.Print(err.Error())
			responseValue.Error = err.Error()
			response = append(response, responseValue)
			gctx.JSON(http.StatusUnprocessableEntity, response)
			return
		}
		if serverID > len(serversData)-1 || serverID < 0 {
			err = fmt.Errorf("error on HTTP-request: \"server_id\" parameter must be in range [0:%d]", len(serversData))
			log.Print(err.Error())
			responseValue.Error = err.Error()
			response = append(response, responseValue)
			gctx.JSON(http.StatusUnprocessableEntity, response)
			return
		}
		responseValue.ID = serverID
		var httpCode int
		if httpCode, ok, err = timeRewindTimepointCheck(serversData[serverID], timepoint); !ok {
			err = fmt.Errorf("%s: %s", errorHeader, err)
			log.Print(err)
			responseValue.Error = err.Error()
			response = append(response, responseValue)
			gctx.JSON(httpCode, response)
			return
		}
		transactionIndex := -1
		for currentIndex, currentTransaction := range History[serversData[serverID].DumpSocketsConfigData.RealSocket].Transactions {
			if timepoint.After(currentTransaction.TransactionTime) {
				transactionIndex = currentIndex
			}
		}
		if transactionIndex == -1 {
			err := fmt.Errorf("error on detected atcual rewind timepoint")
			log.Print(err.Error())
			responseValue.Error = err.Error()
			response = append(response, responseValue)
			gctx.JSON(http.StatusUnprocessableEntity, response)
			return
		}
		responseValue.SettedTimepoint = timepoint.String()
		emulationServers.readWriteMutex.RLock()
		emulationServers.rewindChannels[serverID] <- transactionIndex
		emulationServers.readWriteMutex.RUnlock()
		logString := fmt.Sprintf("Successfully rewinded %d server to %s timepoint", serverID, timepoint.String())
		log.Print(logString)
		response = append(response, responseValue)
		gctx.JSON(http.StatusOK, response)
		return
	}
	for currentIndex, currentServerData := range serversData {
		currentResponse := rewindResponse{
			ID: currentIndex,
		}
		if _, ok, err = timeRewindTimepointCheck(currentServerData, timepoint); !ok {
			err = fmt.Errorf("%s: %s", errorHeader, err)
			log.Print(err)
			currentResponse.Error = err.Error()
			response = append(response, currentResponse)
			continue
		}
		currentTransactionIndex := -1
		for currentIndex, currentTransaction := range History[currentServerData.DumpSocketsConfigData.RealSocket].Transactions {
			if timepoint.After(currentTransaction.TransactionTime) {
				currentTransactionIndex = currentIndex
			}
		}
		if currentTransactionIndex == -1 {
			err = fmt.Errorf("%s: error on detected atcual rewind timepoint", errorHeader)
			log.Print(err)
			currentResponse.Error = err.Error()
			response = append(response, currentResponse)
			continue
		}
		currentResponse.SettedTimepoint = timepoint.String()
		emulationServers.readWriteMutex.RLock()
		emulationServers.rewindChannels[currentIndex] <- currentTransactionIndex
		emulationServers.readWriteMutex.RUnlock()
		logString := fmt.Sprintf("Successfully rewinded %d server to %s timepoint", currentIndex, timepoint.String())
		log.Print(logString)
		response = append(response, currentResponse)
	}
	gctx.JSON(http.StatusOK, response)
}

func getSettingsBuffer() []emulationServerSettings {
	emulationServers.readWriteMutex.RLock()
	serversData := make([]emulationServerSettings, len(emulationServers.serversData))
	copy(serversData, emulationServers.serversData)
	emulationServers.readWriteMutex.RUnlock()
	return serversData
}

func timeRewindTimepointCheck(serverData emulationServerSettings, timepoint time.Time) (httpCode int, result bool, err error) {
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
		httpCode = http.StatusUnprocessableEntity
		return
	}
	result = true
	return
}
