MAKEFLAGS += --no-print-directory

## Determine the OS and set the destination path
OS := $(shell go run scripts/detect_os.go)
ifeq ($(OS), windows)
 BUILD := go build -o aw-sync-agent.exe main.go
 CLEAN := rm -rf $(HOME)/AwSyncAgent

else
 BUILD := go build -o aw-sync-agent main.go
 CLEAN := rm -rf $(HOME)/.config/aw

endif

##General commands
run:
	@go run main.go

build:
	@$(BUILD)

check-os:
	@go run scripts/detect_os.go

clean:
	@$(CLEAN)

format:
	@gofmt -s -w .

##(from linux only) Build executables for both windows and linux
build-all:
	@GOOS=windows GOARCH=amd64 go build -o aw-sync-agent.exe main.go
	@go build -o aw-sync-agent main.go
clean-all:
	@rm -rf aw-sync-agent
	@rm -rf aw-sync-agent.exe



##Service commands
service-install:
	@$(MAKE) build
	@go run main.go -service


### Service commands for linux

service-start:
	@echo "Starting ActivityWatch Sync Agent service..."
	@systemctl --user start aw-sync-agent

service-stop:
	@echo "Stopping ActivityWatch Sync Agent service..."
	@systemctl --user stop aw-sync-agent

service-status:
	@echo "ActivityWatch Sync Agent service status:"
	@systemctl --user  status aw-sync-agent

service-remove:
	@$(MAKE) service-stop
	@echo "Deleting ActivityWatch Sync Agent service..."
	@rm $(HOME)/.config/systemd/user/aw-sync-agent.service
	@systemctl --user daemon-reload
	@$(CLEAN)

service-restart:
	@echo "Restarting ActivityWatch Sync Agent service..."
	@systemctl --user restart aw-sync-agent

service-update:
	@echo "Re-Building ActivityWatch Sync Agent application..."
	@go build -o $(HOME)/.config/aw/aw-sync-agent main.go
	@$(MAKE) service-restart
