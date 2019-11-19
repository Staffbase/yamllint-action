FROM golang:alpine3.8 as builder

RUN apk --update upgrade
RUN apk --no-cache --no-progress add make git gcc musl-dev

WORKDIR /build
COPY . .
RUN go build .


FROM sdesbure/yamllint
RUN pip install --upgrade yamllint
COPY --from=builder /build/yamllint-action /usr/bin/yamllint-action
COPY entrypoint.sh /entrypoint.sh


ENTRYPOINT ["/entrypoint.sh"]
