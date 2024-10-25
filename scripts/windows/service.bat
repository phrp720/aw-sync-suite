@echo off

set APP_NAME=aw-sync-agent
set APP_PATH=C:\aw-sync-agent\agent.exe
set NSSM_PATH=C:\path\to\nssm.exe

echo Building ActivityWatch Sync Agent...
go build -o %APP_PATH%

echo Installing service using NSSM...
%NSSM_PATH% install %APP_NAME% "%APP_PATH%"

echo Configuring service...
%NSSM_PATH% set %APP_NAME% Start SERVICE_AUTO_START
%NSSM_PATH% set %APP_NAME% AppDirectory "C:\aw-sync-agent\agent.exe"
%NSSM_PATH% set %APP_NAME% AppParameters ""

echo Starting service...
%NSSM_PATH% start %APP_NAME%

echo Service %APP_NAME% installed and started successfully.
