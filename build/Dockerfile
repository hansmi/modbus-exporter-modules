FROM docker.io/library/alpine:latest AS alpine

ARG PACKAGE_ARCHIVE

RUN --mount=type=bind,source=${PACKAGE_ARCHIVE},destination=/pkg.tar.gz \
  mkdir /src && \
  tar -C /src --strip-components=1 -xvzf /pkg.tar.gz

FROM scratch

COPY --from=alpine /src/ /

# vim: set sw=2 sts=2 et :
