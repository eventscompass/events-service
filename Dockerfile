# Build
FROM golang:1.21.2-alpine as builder

ENV GO111MODULE=on
ENV GOPRIVATE=github.com/eventscompass
ENV GOOS=linux

WORKDIR /service
COPY . .

RUN CGO_ENABLED=0 go build -o "/tmp/eventsservice" ./src

#####

# Run
FROM scratch
COPY --from=builder /tmp/eventsservice .
EXPOSE 8080
CMD [ "./eventsservice" ]