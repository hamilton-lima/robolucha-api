# prepare builder
FROM golang as builder
COPY . $GOPATH/src/github.com/hamilton-lima/robolucha-api/
WORKDIR $GOPATH/src/github.com/hamilton-lima/robolucha-api/

# get dependancies
RUN go get -d -v

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api
RUN chmod +x /go/bin/api

# start from scratch
FROM scratch
EXPOSE 5000

# Copy our static executable
COPY --from=builder /go/bin/api /go/bin/api
ENTRYPOINT ["/go/bin/api"]
