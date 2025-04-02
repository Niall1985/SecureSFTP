FROM golang:1.22.5-alpine

RUN apk update && apk add --no-cache python3 py3-pip

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# if you have a requirements.txt file, uncomment the next two lines.
#COPY requirements.txt ./
RUN pip3 install --no-cache-dir -r requirements.txt

RUN go build -tags netgo -ldflags '-s -w' -o myapp ./main.go

EXPOSE 8080

CMD ["./myapp"]