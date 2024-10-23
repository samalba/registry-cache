FROM "varnish:7.6-alpine"
USER root
RUN apk add --no-cache --update bash bind-tools
EXPOSE 8080
COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
USER varnish
