FROM gcr.io/distroless/base
COPY kubesearch /usr/bin/kubesearch
ENTRYPOINT ["/usr/bin/kubesearch"]
