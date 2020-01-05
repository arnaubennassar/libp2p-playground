FROM golang:alpine AS build-env
# DOWNLOAD DEPS
ADD go.mod go.sum /src/
WORKDIR /src
RUN go mod download
# BUILD BIN
ADD main.go utils.go /src/
RUN go build -o app
# COPY BIN INTO CLEAN IMAGE
FROM alpine
COPY --from=build-env /src/app /
# RUN BIN. USE -e OPTS="your options here" 
ENTRYPOINT ./app $OPTS
