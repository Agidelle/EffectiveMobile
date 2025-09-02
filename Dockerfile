FROM golang:1.25.0-alpine3.22 AS build

WORKDIR /app

RUN apk add git

ENV APP_NAME="SUBS"

RUN mkdir /out
COPY . /app/

RUN go build  \
    -o /out/${APP_NAME}  \
    github.com/agidelle/effectivemobile


FROM alpine:3.22

WORKDIR /app

ENV SUBS_PORT=3000

COPY --from=build /out/SUBS /app/
COPY --from=build /app/.env /app/
COPY --from=build /app/migrations /app/migrations
COPY --from=build /app/entrypoint.sh /app/

RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]

EXPOSE ${SUBS_PORT}

CMD ["/app/SUBS"]
