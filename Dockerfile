FROM golang:1.23-alpine AS builder
ARG VERSION
ENV VERSION=${VERSION}
RUN apk update && apk add --no-cache git
WORKDIR /dating
COPY . /dating
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.version=${VERSION} -s -w" -o dating cmd/restful/main.go

FROM scratch
COPY --from=builder /dating/dating .
COPY --from=builder /dating/app.env .
EXPOSE 8089
CMD [ "/dating" ]