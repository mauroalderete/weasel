FROM golang:1.19-alpine3.16 as builder

ENV GOPATH=/go

RUN apk add build-base

WORKDIR ${GOPATH}/src/weasel
COPY --chown=root:root ./ .
RUN go generate
RUN go build

FROM golang:1.19-alpine3.16 as runner

ENV GOPATH=/go
ENV APPDIR=/app
ENV GOSRC=${GOPATH}/src
ENV THREAD=2
ENV GATEWAY=https://cloudflare-eth.com

WORKDIR ${APPDIR}/store
WORKDIR ${APPDIR}

COPY --from=builder --chown=root:root ${GOPATH}/src/weasel/weasel /usr/local/bin/weasel

CMD [ "sh", "-c", "weasel run -t ${THREAD} -g ${GATEWAY} --logfile /var/log/weasel/weasel.log --match-verbose --unmatch-verbose --info-verbose --match-file ${APPDIR}/store/match.json" ]
