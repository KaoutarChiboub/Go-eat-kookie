# Use official and latest Golang image
FROM golang:latest

# Set working directory
WORKDIR /app

# Copy the source code of this folder inside the container
COPY . . /app/

# Download and install dependencies
RUN go get -d -v ./...

# Build the GO app
RUN go build -o api .

#Expose the port
EXPOSE 8008

#Run the executable
CMD ["./api"]