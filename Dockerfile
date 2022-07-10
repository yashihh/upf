FROM bitnami/minideb:bullseye

LABEL description="free5GC UPF service" version="Stage 3"

ENV F5GC_MODULE upf
ENV DEBIAN_FRONTEND noninteractive
ARG DEBUG_TOOLS

# Install debug tools ~ 100MB (if DEBUG_TOOLS is set to true)
RUN if [ "$DEBUG_TOOLS" = "true" ] ; then apt-get update && apt-get install -y vim strace net-tools iputils-ping curl netcat tcpdump iptables ; fi

Run useradd free5gc
Run mkdir -p /free5gc && chown -R free5gc:free5gc /free5gc
USER free5gc

# Set working dir
WORKDIR /free5gc
RUN mkdir -p config/ log/

# Copy executable
COPY build/bin/${F5GC_MODULE} ./

# Config files volume
VOLUME [ "/free5gc/config" ]

# Exposed ports
EXPOSE 8000

USER root