FROM golang:alpine AS development
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
ADD . /app
RUN cd /app && go build -o goapp

FROM alpine:latest AS production
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=development /app/goapp /app
EXPOSE 443
ENTRYPOINT ["./goapp"]