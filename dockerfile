FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o company-service

FROM busybox:latest

WORKDIR /company-service/cmd

COPY --from=build /app/cmd/company-service .

COPY --from=build /app/.env /company-service

EXPOSE 8082

CMD ["./company-service"]