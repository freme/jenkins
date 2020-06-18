FROM golang:alpine
RUN apk --no-cache add git
WORKDIR /go/src/github.com/freme/jenkins
RUN git clone https://github.com/freme/jenkins.git .
#RUN GO111MODULE=on go get github.com/urfave/cli/v2
RUN GO111MODULE=on go get -d ./...
WORKDIR /go/src/github.com/freme/jenkins/cmd/getFailedStepsLogs
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o getFailedStepsLogs .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/freme/jenkins/cmd/getFailedStepsLogs/getFailedStepsLogs /bin/getFailedStepsLogs
COPY --from=0 /go/src/github.com/freme/jenkins/docker-entrypoint.sh /bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["getFailedStepsLogs"]

