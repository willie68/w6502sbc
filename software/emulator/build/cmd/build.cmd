@echo off
go build -ldflags="-s -w" -o emu6502.exe cmd/main.go