FROM golang:alpine3.18
LABEL maintainer="Titanio Yudista <titanioyudista98@gmail.com>"
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o fundhub-dashboard-dev .
EXPOSE 8000
CMD ["./fundhub-dashboard-dev"]
