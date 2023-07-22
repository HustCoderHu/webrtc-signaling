package main

import (
	"net/http"
	"webrtc-signaling/pkg/logger"
	"webrtc-signaling/webchan"
)

func main() {
    // logger.InitLog() // linux 上执行这行，log 会写进 syslog，否则就输出到 stdout
    logger.SetLogLevel(logger.LOG_INFO)
    logger.Info("default log level: %d", logger.GetLogLevel())

    server := webchan.NewServer()

    server.OnAuth = func(args interface{}) error {
        logger.Info("OnAuth")
        return nil
    }

    go server.Loop()

    // server.OnConnection = func(c *webchan.Connection, ars interface{}) {
    //     logger.Info("OnConnection from %s", c.PeerInfo())
    // }

    // server.OnMessage = func(c *webchan.Connection, message []byte) {
    //     logger.Info("OnMessage from %s", c.PeerInfo())
    // }

    // server.OnDisconnection = func(c *webchan.Connection, message string) {
    //     logger.Info("OnDisconnection from %s", c.PeerInfo())
    // }

    serveMux := http.NewServeMux()

    serveMux.Handle("/", server)

    port := "8000"
    logger.Info("server listening on: " + port)

    http.ListenAndServe(":" + port, serveMux)
}
