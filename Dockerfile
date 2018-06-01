FROM golang:alpine3.7
RUN apk --no-cache add git
WORKDIR /go/src/github.com/freme/jenkins
RUN git clone https://github.com/freme/jenkins.git .
RUN go get -d ./...
WORKDIR /go/src/github.com/freme/jenkins/cmd/getFailedStepsLogs
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o getFailedStepsLogs .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/freme/jenkins/cmd/getFailedStepsLogs/getFailedStepsLogs /bin/getFailedStepsLogs
COPY docker-entrypoint.sh /bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["getFailedStepsLogs"]

