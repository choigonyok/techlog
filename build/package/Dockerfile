FROM --platform=amd64 golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN cd ../../cmd/main.go && go build main.go

CMD [ "../../cmd/main" ]
