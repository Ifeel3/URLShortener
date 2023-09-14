FROM golang:alpine
ADD src/ /server/
WORKDIR /server/
RUN go get github.com/redis/go-redis/v9
RUN go build -o shortener .
CMD ./shortener