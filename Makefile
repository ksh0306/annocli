build:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o ${GOPATH}/bin/annocli-win-amd64.exe .
	GOOS=linux GOARCH=amd64 go build -o ${GOPATH}/bin/annocli-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o ${GOPATH}/bin/annocli-darwin-amd64 .
	
reset:
	rm $(HOME)/annocli/config.yaml
