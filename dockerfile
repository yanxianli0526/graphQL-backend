# build stage
FROM golang:1.17.6-alpine3.15 AS graphql-go-template

ARG CI_JOB_TOKEN

ADD . /src
RUN apk add --no-cache git && \
    git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.smart-aging.tech/".insteadOf "https://gitlab.smart-aging.tech/" && \
    cd /src && go build -o server 

# final stage
FROM alpine:3.15
WORKDIR /app
COPY --from=graphql-go-template /src/server /app/server
COPY ./zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip 
EXPOSE 4000
CMD ["/app/server"]