FROM golang:1.14 as builder

ENV GRPC_HEALTH_PROBE_VERSION=v0.3.3

RUN \
  wget -qO/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
  && chmod +x /bin/grpc_health_probe

FROM scratch

ARG VERSION

LABEL maintainer="vontikov@gmail.com"
LABEL version="${VERSION}"

COPY --from=builder /bin/grpc_health_probe /usr/local/bin/
COPY stoa /usr/local/bin/
CMD ["stoa"]
