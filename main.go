package main

import (
	"net/http"
	"pass_web/internal/api/router"
)

func main() {
    err := http.ListenAndServe(":8080", router.NewMutexHandler())
    if err != nil{
        panic(err)
    }
}
