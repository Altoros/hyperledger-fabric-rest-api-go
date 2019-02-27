# Stage 1 - building app

FROM golang:1.12-alpine3.9 AS build
RUN go version

COPY . /src/
WORKDIR /src/

RUN apk add --no-cache make git

RUN make build

# Stage 2 - serving app

FROM alpine:3.9

RUN adduser -S -D -H -h /app appuser
USER appuser

WORKDIR /app/
COPY --from=build /src/configs/ configs/
COPY --from=build /src/test/ test/
COPY --from=build /src/build/frag .

EXPOSE 8080

ENTRYPOINT ["/app/frag"]