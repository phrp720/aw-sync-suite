MAKEFLAGS += --no-print-directory

## Determine the OS and set the destination path
OS := $(shell go run scripts/detect_os.go)
ifeq ($(OS), windows)
 BUILD := GOOS=windows go build -o C:/aw-sync-agent/agent.exe main.go && cp -r .env C:/aw-sync-agent/.env
 DEV_BUILD := GOOS=windows go build -o bin/agent.exe main.go
 CLEAN := rm -rf C:/aw-sync-agent
 ##Service commands
 SERVICE := scripts/windows/service.bat
else
 BUILD := sudo go build -o /bin/agent main.go && sudo cp -r .env /bin/.env
 DEV_BUILD := go build -o bin/agent main.go
 CLEAN :=  sudo rm -rf /bin/agent && sudo rm -rf /bin/.env
 ##Service commands
 SERVICE := sudo scripts/linux/service.sh
endif

##Development commands
dev-run:
	@go run main.go -asService=false -awUrl=http://localhost:5600 -prometheusUrl=http://localhost:9090 -userID=DevUser

dev-build:
	@$(DEV_BUILD)
	@cp -r .env bin/.env

dev-clean:
	@rm -rf bin


##General commands
run:
	@go run main.go

build:
	@$(BUILD)


check-os:
	@go run scripts/detect_os.go

test:
	@go test -v ./...

clean:
	@$(CLEAN)

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
	@echo "Deleting ActivityWatch Sync Agent service..."
	@sudo rm /etc/systemd/system/aw-sync-agent.service
	@sudo systemctl daemon-reload

service-restart:
	@echo "Restarting ActivityWatch Sync Agent service..."
	@sudo systemctl restart aw-sync-agent

service-update:
	@echo "Re-Building ActivityWatch Sync Agent application..."
	@go build -o $APP_PATH main.go
	@sudo cp -r .env /bin/.env
	@$(MAKE) service-restart
