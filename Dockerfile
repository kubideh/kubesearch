# Dockerfile is used to build a container image for the kubesearch
# server process. The container image is based on "distroless".
FROM gcr.io/distroless/base
COPY kubesearch /usr/bin/kubesearch
ENTRYPOINT ["/usr/bin/kubesearch"]
