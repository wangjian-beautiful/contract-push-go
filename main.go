package main

import (
	"gcex-contract-go/marker"
	"gcex-contract-go/server"
	"github.com/gin-gonic/gin"
	_ "strings"
)

func main() {
	nominalServer := server.NewNominalServer()
	authServer := server.NewAuthServer()
	nominalServer.Start()
	authServer.Start()

	marker.AddWSService(nominalServer)
	marker.AddWSService(authServer)
	marker.StarterMqPush()
	marker.StarterDepth()
	marker.StarterFundingRatePush()
	marker.StarterPositionPush()

	r := gin.Default()
	r.GET("/ws", nominalServer.Ws.Handler())
	r.GET("/auth", authServer.Ws.Handler())
	r.Run(":9090")

}
