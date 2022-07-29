FROM golang:1.18

WORKDIR /usr/src/app


##  REQUIRS
##  TTOKEN    - telegram bot token
##  BITRIX_TOKEN  - token to acces bitrix api
##  ADMIN_ID   -  admin's telegram chat id

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/offset-tg-bot .

CMD ["offset-tg-bot"]