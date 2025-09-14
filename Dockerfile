FROM golang:1.21-alpine


WORKDIR /app
COPY . .


RUN go mod download


CMD ["go", "run", "./cmd/server"]