# Stage 1 - building app

FROM golang:1.12-alpine3.9 AS build

RUN apk add --no-cache make git

WORKDIR /src/

# download modules in separated layer, to speed up rebuild by utilising Docker layer caching system
COPY go.mod .
COPY go.sum .
# NOTE: build error may occur due to temporary unavailability of some packages sources
# Wait and build again is usually a good solution
RUN go mod download

COPY . /src/
RUN make build


# Stage 2 - serving app

FROM alpine:3.9

WORKDIR /app/

COPY --from=build /src/build/frag .

EXPOSE 8080

ENTRYPOINT ["/app/frag"]