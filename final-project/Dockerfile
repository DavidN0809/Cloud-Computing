FROM golang:1.21-alpine AS build
WORKDIR /src/
COPY taskmanagement.go /src/
RUN go mod init github.com/DavidN0809/Cloud-Computing/final-project
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/taskmanagement
FROM scratch
COPY --from=build /bin/taskmanagement /bin/taskmanagement
ENTRYPOINT ["/bin/taskmanagement"]
