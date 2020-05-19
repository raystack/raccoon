package app

import (
	"clickstream-service/env"
	"clickstream-service/router"
	"fmt"
	"github.com/urfave/negroni"
	"log"
	"strconv"
)

func StartHTTPServer() {

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(router.Router())
	port := fmt.Sprintf(":%s", strconv.Itoa(env.AppPort()))

	log.Println("Starting server on port", port)
	go server.Run(port)
}
