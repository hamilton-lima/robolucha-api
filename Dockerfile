# prepare builder
FROM golang:1.12.6 as builder

RUN mkdir -p /usr/local/share/robolucha-api
WORKDIR /usr/local/share/robolucha-api/

# get dependencies
COPY go.mod /usr/local/share/robolucha-api/
RUN go get 

# copy source code
COPY . /usr/local/share/robolucha-api/

# copy gamedefinition files
RUN cp -r gamedefinition /tmp/gamedefinition

# copy grade files
RUN cp -r grade /tmp/grade

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/api
RUN chmod +x /tmp/api

# start from scratch
FROM alpine
EXPOSE 5000

RUN mkdir -pv /usr/src/app/gamedefinition
RUN mkdir -pv /usr/src/app/grade

# Copy our static executable
COPY --from=builder /tmp/api /usr/src/app
COPY --from=builder /tmp/gamedefinition /usr/src/app/gamedefinition
COPY --from=builder /tmp/grade /usr/src/app/grade
RUN ls -alhR /usr/src/app
ENTRYPOINT ["/usr/src/app/api", "/usr/src/app/gamedefinition", "/usr/src/app/grade"]
