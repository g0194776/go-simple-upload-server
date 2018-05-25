FROM alpine:3.5

COPY go-simple-upload-server /usr/local/bin/app
RUN chmod +x /usr/local/bin/app