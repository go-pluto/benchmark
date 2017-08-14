FROM alpine:3.6

RUN mkdir /config /users
RUN apk add -U ca-certificates

ADD ./benchmark /bin/

ENTRYPOINT ["/bin/benchmark"]
CMD ["-config", "/config/config.toml", "-userdb", "/users/passwd", "-logtostderr=true", "-v=0"]
