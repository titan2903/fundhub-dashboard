FROM golang:alpine3.18
LABEL maintainer="Titanio Yudista <titanioyudista98@gmail.com>"
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY .env .
COPY . .
RUN go build -o dashboard .
EXPOSE 8000
CMD ["./dashboard"]
