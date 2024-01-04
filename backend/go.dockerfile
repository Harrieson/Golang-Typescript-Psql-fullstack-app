FROM golang:1.20
WORKDIR /app
COPY . .

# Download and install deps 
RUN go get -d -v ./...

#Build the App 

RUN go build -o api .

EXPOSE 8000

CMD ["./api"]