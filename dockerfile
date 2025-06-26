# Build stage
FROM golang:1.24-bookworm AS builder

ENV DEBIAN_FRONTEND=noninteractive
ENV NET_SNMP_VERSION=5.9.4

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    make \
    libssl-dev \
    libperl-dev \
    python3-dev \
    python3-setuptools \
    autoconf \
    automake \
    libtool \
    pkg-config \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Download and build net-snmp
WORKDIR /tmp
RUN wget https://sourceforge.net/projects/net-snmp/files/net-snmp/${NET_SNMP_VERSION}/net-snmp-${NET_SNMP_VERSION}.tar.gz/download -O net-snmp-${NET_SNMP_VERSION}.tar.gz \
    && tar -xzf net-snmp-${NET_SNMP_VERSION}.tar.gz \
    && cd net-snmp-${NET_SNMP_VERSION} \
    && ./configure \
        --prefix=/usr/local \
        --enable-shared \
        --disable-static \
        --with-default-snmp-version="3" \
        --with-sys-contact="root@127.0.0.1" \
        --with-sys-location="Unknown" \
        --with-logfile="/var/log/snmpd.log" \
        --with-persistent-directory="/var/net-snmp" \
        --disable-embedded-perl \
        --without-perl-modules \
        --disable-manuals \
        --with-openssl \
        --disable-mib-loading \
        --with-mib-modules="mibII/system_mib,snmpv3mibs,ucd_snmp" \
    && make -j$(nproc) \
    && make install

# Runtime stage
FROM golang:1.24-bookworm

ENV DEBIAN_FRONTEND=noninteractive
ENV GO_TEST_TIMEOUT=120s
ENV PATH="/usr/local/bin:${PATH}"


RUN apt-get update && apt-get install -y \
    libssl3 \
    git \
    && rm -rf /var/lib/apt/lists/*


COPY --from=builder /usr/local /usr/local


RUN ldconfig

RUN mkdir -p /var/run \
    && mkdir -p /var/net-snmp \
    && mkdir -p /var/log


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Command to run tests
CMD ["go", "test", "-v", "./..."]