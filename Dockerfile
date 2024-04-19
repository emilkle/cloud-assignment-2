FROM golang:1.20

LABEL authors="Marius Robsham Dahl, Kristian Welle Glomset, Emil Klevgård-Slåttsveen"

WORKDIR /go/src/app

COPY ./database /go/src/app/database
COPY ./firestoreEmulator /go/src/app/firestoreEmulator
COPY ./functions /go/src/app/functions
COPY ./handlers /go/src/app/handlers
COPY ./resources /go/src/app/resources
COPY ./tests /go/src/app/tests
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
COPY ./main.go /go/src/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main

EXPOSE 8080
EXPOSE 8081

CMD ["./main"]