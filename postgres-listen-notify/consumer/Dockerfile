FROM golang:1.21 as build

COPY . /src

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o consumer .

FROM scratch as production

COPY --from=build /src/consumer .

CMD [ "/consumer" ]