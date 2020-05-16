FROM alpine:3.11

RUN apk update \
	&& apk -U upgrade \
	&& apk add --no-cache ca-certificates \
	&& update-ca-certificates --fresh \
	&& rm -rf /var/cache/apk/*

# adds app user and fix app folder's permission
RUN	addgroup app \
	&& adduser -S app -u 1000 -G app

USER app

COPY --chown=app:app app /usr/local/bin/
RUN chmod +x /usr/local/bin/app

ENTRYPOINT [ "/usr/local/bin/app" ]
