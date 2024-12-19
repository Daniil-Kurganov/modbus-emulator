package src

import (
	"github.com/gin-gonic/gin"
)

func StartHTTPServer() {
	router := gin.Default()
	emulator := router.Group("/modbus-emulator")
	{
		emulator.GET("servers") // id of working emulations servers
		time := emulator.Group("time")
		{
			time.GET("actual")            // info about current time of all working emulations servers
			time.GET("start&end")         // info about starting and ending times all working emulations servers
			time.POST("rewind")           // rewind all working servers on a specified time
			time.GET(":server/actual")    // info of current time of server
			time.GET(":server/start&end") // info about starting and ending of server
			time.POST(":server/rewind")   // rewind server on a specified time
		}
		settings := emulator.Group("settings") // emulation socket, dump socket, protocol, emulation mode (one-time or continuously)
		{
			settings.GET("")                        // settings for all working servers
			settings.POST("emulation_mode")         // set emulation mode for all server
			settings.GET(":server")                 // settings of server
			settings.POST(":server/emulation_mode") // set emulation mode for server
		}
	}
	router.Run("127.0.0.1:8080")
}
