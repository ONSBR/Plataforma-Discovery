default:
		GOOS=linux CGO_ENABLED=0 go build -o dist/discovery

convey:
	goconvey --port 8890

test:
	go test ../... -v



