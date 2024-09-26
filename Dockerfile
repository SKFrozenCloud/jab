FROM golang:latest

COPY . .

RUN chmod +x ./externals/upx-4.2.4
RUN make build

CMD ["./main"]
