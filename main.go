package main

import "HiChat/initialize"

func main() {
	initialize.InitDB()
	initialize.InitLogger()
}
