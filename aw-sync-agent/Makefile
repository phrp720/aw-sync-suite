MAKEFLAGS += --no-print-directory

## Determine the OS and set the destination path
OS := $(shell go run scripts/detect_os.go)
ifeq ($(OS), windows)
 BUILD := go build -o aw-sync-agent.exe main.go
 CLEAN := rm -rf aw-sync-agent.exe
 CLEAN-SERVICE := rm -rf $(HOME)/AwSyncAgent

 SERVICE-START := $(HOME)/AwSyncAgent/nssm.exe start aw-sync-agent
 SERVICE-STOP := $(HOME)/AwSyncAgent/nssm.exe stop aw-sync-agent
 SERVICE-REMOVE := $(HOME)/AwSyncAgent/nssm.exe remove aw-sync-agent confirm && $(CLEAN-SERVICE)
 SERVICE-RESTART := $(HOME)/AwSyncAgent/nssm.exe restart aw-sync-agent
 SERVICE-STATUS := $(HOME)/AwSyncAgent/nssm.exe status aw-sync-agent

else
 BUILD := go build -o aw-sync-agent main.go
 CLEAN := rm -rf aw-sync-agent
 CLEAN-SERVICE := rm -rf $HOME/.config/aw

 SERVICE-START := systemctl --user start aw-sync-agent
 SERVICE-STOP := systemctl --user stop aw-sync-agent
 SERVICE-REMOVE := rm $HOME/.config/systemd/user/aw-sync-agent.service && systemctl --user daemon-reload
 SERVICE-RESTART := systemctl --user restart aw-sync-agent
 SERVICE-STATUS := systemctl --user status aw-sync-agent
endif

##General commands
run:
	@go run main.go

build:
	@$(BUILD)

clean:
	@$(CLEAN)

check-os:
	@go run scripts/detect_os.go

clean-service:
	@echo "Cleaning ActivityWatch Sync Agent service files..."
	@$(CLEAN-SERVICE)

format:
	@gofmt -s -w .

test:
	@go test ./tests/

##(from linux only) Build executables for both windows and linux
build-all:
	@GOOS=windows GOARCH=amd64 go build -o aw-sync-agent.exe main.go
	@go build -o aw-sync-agent main.go

clean-all:
	@rm -rf aw-sync-agent
	@rm -rf aw-sync-agent.exe



##Service commands
service-install: build
	@go run main.go -service


### Service commands for linux

service-start:
	@echo "Starting ActivityWatch Sync Agent service..."
	@$(SERVICE-START)

service-stop:
	@echo "Stopping ActivityWatch Sync Agent service..."
	@$(SERVICE-STOP)

service-status:
	@echo "ActivityWatch Sync Agent service status:"
	@$(SERVICE-STATUS)

service-remove: service-stop
	@echo "Deleting ActivityWatch Sync Agent service..."
	@$(SERVICE-REMOVE)
	@$(CLEAN-SERVICE)

service-restart:
	@echo "Restarting ActivityWatch Sync Agent service..."
	@$(SERVICE-RESTART)

