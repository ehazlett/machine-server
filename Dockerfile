FROM scratch
COPY machine-server /machine-server
COPY swagger /swagger
WORKDIR /
EXPOSE 8080
ENTRYPOINT ["/machine-server"]
