FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy go files in the ROOT DIR and the three modules
COPY *.go ./
COPY scraper/ ./scraper/
COPY indexer/ ./indexer/
COPY server/ ./server/

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o scrape-psgtech

EXPOSE 8000

CMD ["./scrape-psgtech"]
