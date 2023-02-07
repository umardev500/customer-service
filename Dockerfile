FROM golang as dev

WORKDIR /app

COPY . .

EXPOSE 5012

CMD air
