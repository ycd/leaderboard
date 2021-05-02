FROM golang:alpine

WORKDIR /src
COPY go.mod go.sum ./

RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app cmd/main.go

FROM alpine
COPY --from=0 /bin/app /bin/app 

ENV PORT=8000
EXPOSE 8000
ENTRYPOINT [ "/bin/app" ]