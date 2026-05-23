FROM docker.io/golang:1.26.2-alpine AS build

WORKDIR /src/
RUN apk add git
COPY go.* .
RUN go mod download
COPY enemies enemies
COPY logging logging
COPY redis /redis
COPY *.go .
RUN go build -v -o nullscaple

FROM docker.io/alpine
RUN apk add --no-cache tzdata
ENV TZ=America/New_York
RUN cp /usr/share/zoneinfo/America/New_York /etc/localtime
COPY enemies.json enemies.json
COPY static /static
COPY templates /templates

COPY --from=build /src/nullscaple /nullscaple

ENTRYPOINT [ "/nullscaple" ]