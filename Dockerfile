FROM golang:1.26-alpine
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main main.go

EXPOSE 8090
CMD [ "./main" ]
