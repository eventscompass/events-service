# This Dockerfile uses a multi-stage build.
# https://docs.docker.com/build/building/multi-stage/
#
# In order to keep the production image as small as possible we
# will have two separate stages: one for building the binary, and
# another for running that binary.
# The first stage uses a full-scale Go image to build the binary.
# The second stage uses the scratch image as its base and copies
# just the built binary from the previous stage. None of the
# build tools required to build the application are included in
# the resulting image.


#---------------------------------------------------------------#
# Build
#
# For compiling the service binary we will use the latest golang
# image. Name the build stage so that we can refer to it from
# later stages.
FROM golang:1.21.2-alpine AS builder

# Set the working directory inside the docker container.
# Note that if we are not careful we might accidentally overwrite
# some files or folders of the default folder. The `/usr` folder
# is considered a safe place to put the service files.
# The folder `/usr/service` will be created if it does not exists.
WORKDIR /usr/service

# Install all the dependencies inside the docker image workdir.
# Note that the dependencies are committed to the repo, so we can
# simply copy them instead of downloading with `go mod download`.
COPY go.mod go.sum ./
COPY ./vendor ./vendor/

# Copy the src code.
# We could also do simply `COPY ./ ./`, but splitting the copy
# operation into two different steps improves cache reuse.
# Changing the src files during development and re-building the
# image will not invalidate the entire copy step.
COPY ./src ./src/

# Compile the service binary.
RUN \
    CGO_ENABLED=0 \
    GOOS=linux \
    go build -o "./eventsservice" ./src


#---------------------------------------------------------------#
# Run
#
# Run the service binary with the scratch image as base.
FROM scratch

# Copy just the built artifact from the previous stage into this
# new stage. Any intermediate artifacts are left behind.
COPY --from=builder /usr/service/eventsservice ./

# The `EXPOSE` function does not publish the port.
# But we can document in the Dockerfile what ports the
# application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#
# The REST server listens on port 8080, and the gRPC - on 8081.
EXPOSE 8080/tcp
EXPOSE 8081/tcp

# Run the service binary.
CMD [ "./eventsservice" ]
