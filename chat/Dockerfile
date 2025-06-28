FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . /app

RUN go mod download \
    && go get github.com/johseongeon/chat_package \
    && go build -o main .

WORKDIR /dist

RUN cp /app/main .

FROM scratch

COPY --from=builder /dist/main .

EXPOSE 8080

ENV MONGO_URI=mongodb://host.docker.internal:27017

ENTRYPOINT ["/main"]