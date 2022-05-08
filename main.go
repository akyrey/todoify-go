package main

import (
	"fmt"
	"log"
	"net/http"
    "github.com/akyrey/todoify-go/handlers"
)

func main() {
    server := &http.Server{
        Addr: fmt.Sprint(":8000"),
        Handler: handlers.New(),
    }

    log.Printf("Starting http server. Listening on %q", server.Addr)

    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Printf("%v", err)
    } else {
        log.Println("Server closed!")
    }
}
