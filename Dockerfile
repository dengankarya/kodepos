FROM golang:1.23-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /kodepos .


FROM scratch

COPY --from=build /kodepos /kodepos

EXPOSE 3000

ENTRYPOINT ["/kodepos"]
