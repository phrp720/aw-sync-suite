package service

import "log"

func CreateLinuxService() {
	//Here we will handle the windows and linux service creation
	//We will use the nssm for windows and the systemd for linux
	//We will create a service that will run the agent as a service and with -service flag we will pass all the data to the excutable
	log.Print("Running as a service In Linux...")

}
