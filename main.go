package main

import (
	"fmt"
	"net/http"
	"pass_web/internal/api/router"
)

func main() {
	fmt.Println("YOP")

    http.FileServer(http.Dir("static"))
    err := http.ListenAndServe("localhost:8080", router.NewMutexHandler())
    if err != nil{
        panic(err)
    }
}
