package service

import "log"

func CreateWindowsService() {
	//Here we will handle the windows and linux service creation
	//We will use the nssm for windows and the systemd for linux
	//We will create a service that will run the agent as a service and with -service flag we will pass all the data to the excutable
	//os.Exit(0)
	log.Print("Running as a service In Windows...")

}
