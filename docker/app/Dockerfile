FROM likechain/base as builder

WORKDIR /likechain
COPY . .
RUN go build -o /go/bin/liked cmd/liked/main.go
RUN go build -o /go/bin/likecli cmd/likecli/main.go

FROM alpine:latest

ARG UID=1000
ARG GID=1000

RUN apk add ca-certificates
RUN addgroup -g $GID likechain && \
    adduser --uid $UID --disabled-password --gecos "" --ingroup likechain --home /likechain likechain
USER $UID
WORKDIR /bin
COPY --from=builder /go/bin/liked .
COPY --from=builder /go/bin/likecli .
