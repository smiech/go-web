FROM golang:1.6

RUN go get github.com/adam72m/go-web

# Expose the application on port 8080
EXPOSE 8080

# Set the entry point of the container to the bee command that runs the
# application and watches for changes
CMD ["go-web", "run"]