FROM golang:1.23rc1-alpine3.20

WORKDIR /movie-rating-api-go

COPY . /movie-rating-api-go/

RUN go mod tidy

EXPOSE 4040

CMD [ "go", "run", "./cmd/main.go" ]