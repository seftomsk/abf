FROM golang:1.17.3 as build

ENV BIN_FILE /opt/abf/abf-app
ENV CODE_DIR /go/src

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ./cmd/abf/*

FROM alpine:3.15.0

ENV BIN_FILE "/opt/abf/abf-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}
CMD ${BIN_FILE}
