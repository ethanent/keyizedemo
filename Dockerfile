FROM golang:1.15.6-alpine

WORKDIR /buildarea
COPY . .

RUN ["go", "mod", "download"]
RUN ["go", "build", "-o", "out", "."]

FROM alpine:latest
RUN ["apk", "--no-cache", "add", "ca-certificates"]

WORKDIR /main
COPY --from=0 /buildarea/out .
COPY --from=0 /buildarea/static ./static/

CMD ["./out"]
