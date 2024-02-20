FROM oven/bun:alpine
COPY service ./service
RUN apk update && \
    apk add --no-cache runit
RUN rm -rf /var/cache/apk/*
EXPOSE 3000
EXPOSE 3001
ENTRYPOINT ["runsvdir", "./service"]