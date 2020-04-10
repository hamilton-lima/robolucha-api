# prepare builder
FROM golang:1.12.6 as builder

RUN mkdir -p /usr/local/share/robolucha-api
WORKDIR /usr/local/share/robolucha-api/

# get dependencies
COPY go.mod /usr/local/share/robolucha-api/
RUN go get 

# copy source code
COPY . /usr/local/share/robolucha-api/

# copy metadata files
RUN cp -r metadata /tmp/metadata

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/api
RUN chmod +x /tmp/api

# start from scratch
FROM alpine
EXPOSE 5000

RUN mkdir -pv /usr/src/app/metadata

# Copy our static executable
COPY --from=builder /tmp/api /usr/src/app
COPY --from=builder /tmp/metadata /usr/src/app/metadata
RUN ls -alhR /usr/src/app
ENTRYPOINT ["/usr/src/app/api", "/usr/src/app/metadata"]
