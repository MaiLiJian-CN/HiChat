package main

import (
	"HiChat/initialize"
	"HiChat/router"
)

func main() {
	initialize.InitDB()
	initialize.InitLogger()

	router:=router.Router()
	router.Run(":8089")
}
