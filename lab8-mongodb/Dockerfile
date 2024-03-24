FROM golang:1.21-alpine AS build
WORKDIR /src/
COPY mongo.go /src/
RUN go mod init github.com/hburnt/CloudNativeCourse
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/mongo-server
FROM scratch
COPY --from=build /bin/mongo-server /bin/mongo-server
ENTRYPOINT ["/bin/mongo-server"]
