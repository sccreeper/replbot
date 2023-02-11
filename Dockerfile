ARG bot_token

FROM alpine
RUN apk add --no-cache go

COPY . .

RUN go build -o discordreplbot

CMD ["./discordreplbot"]