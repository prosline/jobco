FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/prosline/jobco

ENV USER marciodasilva
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET O05cjw1hYoe3VSkX

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://marciodasilva@localhost:5432/jobco?sslmode=disable

WORKDIR /go/src/github.com/prosline/jobco

RUN godep go build

EXPOSE 8888
CMD ./jobco