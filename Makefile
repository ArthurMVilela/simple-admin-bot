include .env

run-bot:
	go run cmd/bot/main.go --token=$(TOKEN)