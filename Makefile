MAKEFLAGS += --no-print-directory

## Determine the OS and set the destination path
OS := $(shell go run scripts/detect_os.go)
ifeq ($(OS), windows)
 BUILD := go build -o aw-sync-agent.exe main.go
 CLEAN := rm -rf C:/aw-sync-agent
 ##Service commands
 SERVICE := scripts/windows/service.bat
else
 BUILD := go build -o aw-sync-agent main.go
 CLEAN :=  sudo rm -rf /opt/aw
 ##Service commands
 SERVICE :=go run main.go -service
endif

##General commands
run:
	@go run main.go

build:
	@$(BUILD)

check-os:
	@go run scripts/detect_os.go

clean:
	@rm -rf aw-sync-agent

format:
	@gofmt -s -w .

##(from linux only) Build executables for both windows and linux
build-all:
	@go build -o aw-sync-agent.exe main.go GOOS=windows GOARCH=amd64
	@go build -o aw-sync-agent main.go
clean-all:
	@rm -rf aw-sync-agent
	@rm -rf aw-sync-agent.exe



##Service commands
service-install:
	@$(MAKE) build
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
	@sudo go build -o /opt/aw/aw-sync-agent main.go
	@$(MAKE) service-restart
