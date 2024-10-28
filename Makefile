MAKEFLAGS += --no-print-directory

## Determine the OS and set the destination path
OS := $(shell go run scripts/detect_os.go)
ifeq ($(OS), windows)
 BUILD := GOOS=windows go build -o aw-sync-agent.exe main.go
 CLEAN := rm -rf C:/aw-sync-agent
 ##Service commands
 SERVICE := scripts/windows/service.bat
else
 BUILD := go build -o aw-sync-agent main.go
 CLEAN :=  sudo rm -rf /bin/aw
 ##Service commands
 SERVICE := sudo scripts/linux/service.sh
endif


##General commands
run:
	@go run main.go

build:
	@$(BUILD)

check-os:
	@go run scripts/detect_os.go

clean:
	@rm -rf agent*

format:
	@gofmt -s -w .

##Service commands
service:
	@$(SERVICE)

service-start:
	@echo "Starting ActivityWatch Sync Agent service..."
	@sudo systemctl start aw-sync-agent

service-stop:
	@echo "Stopping ActivityWatch Sync Agent service..."
	@sudo systemctl stop aw-sync-agent

service-status:
	@echo "ActivityWatch Sync Agent service status:"
	@sudo systemctl status aw-sync-agent

service-remove:
	@$(MAKE) service-stop
	@echo "Deleting ActivityWatch Sync Agent service..."
	@sudo rm /etc/systemd/system/aw-sync-agent.service
	@sudo systemctl daemon-reload
	@$(CLEAN)

service-restart:
	@echo "Restarting ActivityWatch Sync Agent service..."
	@sudo systemctl restart aw-sync-agent

service-update:
	@echo "Re-Building ActivityWatch Sync Agent application..."
	@sudo go build -o /bin/aw/agent main.go
	@$(MAKE) service-restart
