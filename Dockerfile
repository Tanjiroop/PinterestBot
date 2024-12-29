FROM golang:1.22-bullseye

RUN mkdir /bot
WORKDIR /bot
RUN cd /bot
COPY . .
RUN go build .

CMD ["./PinterestBot"]
