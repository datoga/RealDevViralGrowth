FROM golang AS builder

ENV GO111MODULE=on

RUN mkdir -p /go/src/ViralGrowth
WORKDIR /go/src/ViralGrowth

ADD . /go/src/ViralGrowth

WORKDIR /go/src/ViralGrowth

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s -extldflags "-static"' .

FROM alpine

RUN mkdir /app 
WORKDIR /app
COPY --from=builder /go/src/ViralGrowth/RealDevViralGrowth .

CMD ["/app/RealDevViralGrowth"]
