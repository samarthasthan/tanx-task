FROM golang:latest

WORKDIR /usr/src/tanx

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN migrate --version
RUN apt-get install -y make
RUN go build -v -o app ./cmd/tanx/