# prepare builder
FROM golang as builder
COPY go-get.sh $GOPATH/src/gitlab.com/robolucha/robolucha-api/
WORKDIR $GOPATH/src/gitlab.com/robolucha/robolucha-api/

# Authorize SSH Host
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan gitlab.com > /root/.ssh/known_hosts

# Add the keys and set permissions
COPY .ssh/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa 

# Set git to use SSH instead of HTTPS
RUN git config --global --add url."git@gitlab.com:robolucha/robolucha-api.git".insteadOf "https://gitlab.com/robolucha/robolucha-api.git"

# get dependencies
RUN bash -c $GOPATH/src/gitlab.com/robolucha/robolucha-api/go-get.sh

COPY . $GOPATH/src/gitlab.com/robolucha/robolucha-api/

# copy gamedefinition files
RUN cp -r gamedefinition /tmp/gamedefinition

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api
RUN chmod +x /go/bin/api

# start from scratch
FROM alpine
EXPOSE 5000

RUN mkdir -pv /usr/src/app/gamedefinition

# Copy our static executable
COPY --from=builder /go/bin/api /usr/src/app
COPY --from=builder /tmp/gamedefinition /usr/src/app/gamedefinition
RUN ls -alh /usr/src/app
ENTRYPOINT ["/usr/src/app/api", "/usr/src/app/gamedefinition"]
