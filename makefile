updatego:
	go get && go get -u ./... && go mod tidy && go get

build:
	go build -ldflags="-w -s" -o ./main . && ./externals/upx-4.2.4 -9 ./main