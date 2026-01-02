FROM golang AS builder
WORKDIR /src

COPY go.mod ./
COPY ./*.go ./

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app .

FROM scratch
WORKDIR /

COPY --from=builder /out/app /app

USER 1000:1000

EXPOSE 8080
CMD ["/app"]