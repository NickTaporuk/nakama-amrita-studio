DOCKER_IMAGE_NAME=nakama-amrita-studio
TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI1MjIyNzVhOC1hZTQyLTRiM2MtYTdiMS0wOWViMDEyY2Y0MzYiLCJ1c24iOiJteWN1c3RvbXVzZXJuYW1lIiwiZXhwIjoxNTk0MzAzNTMyfQ.3gzXoSmFCeXC1yMs4yVW0M6erabd2sKzP0rzs47FFsY
build-plugin:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin user.plugin.go

# should be run from the root directory
build-image:
	docker build -t $(DOCKER_IMAGE_NAME):latest --no-cache -f ./Dockerfile .

create-user:
	curl "http://127.0.0.1:7350/v2/account/authenticate/device?create=true&username=admin" --user 'defaultkey:' --data '{"id":"uniqueidentifier"}'
req:
	curl "http://127.0.0.1:7350/v2/rpc/add_user" -H 'authorization: Bearer $(TOKEN)' -d '"{\"tournament_id\": \"34fac801-75fd-4f01-ac6f-d2e7254972e3\"}"'