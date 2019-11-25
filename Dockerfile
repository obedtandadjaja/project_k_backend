# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.15.0 as builder

ENV GO111MODULE=on
RUN mkdir -p $GOPATH/src/github.com/obedtandadjaja
WORKDIR $GOPATH/src/github.com/obedtandadjaja

ADD . .
RUN go mod download
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
CMD /bin/app migrate; /bin/app
