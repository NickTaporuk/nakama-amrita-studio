DOCKER_IMAGE_NAME=nakama-amrita-studio
TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI1MjIyNzVhOC1hZTQyLTRiM2MtYTdiMS0wOWViMDEyY2Y0MzYiLCJ1c24iOiJteWN1c3RvbXVzZXJuYW1lIiwiZXhwIjoxNTk0MzExMDgzfQ.xfdpnmIdWB82JlwGso9kVqS0HJnNQPgEB8kg6zpPLu4
build-plugin:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin user.plugin.go

# should be run from the root directory
build-image:
	docker build -t $(DOCKER_IMAGE_NAME):latest --no-cache -f ./Dockerfile .

# get token to work with add user to a tournament
create-user:
	curl "http://127.0.0.1:7350/v2/account/authenticate/device?create=true&username=admin" --user 'defaultkey:' --data '{"id":"uniqueidentifier"}'

# add user to a tournament
add-user-to-tournament:
	curl "http://127.0.0.1:7350/v2/rpc/addUserToTournament" -H 'authorization: Bearer $(TOKEN)' -d '"{\"tournament_id\": \"34fac801-75fd-4f01-ac6f-d2e7254972e1\"}"'