#!/bin/bash
set -e

VERSION=v0.1-alpha
OS=darwin
ARCH=amd64

BINARY_NAME=rubberdoc
BINARY_DST=$(eval echo \$\{GOPATH\})/bin/

PACKAGE_NAME=${BINARY_NAME}-${VERSION}.${OS}-${ARCH}.tar.gz
PACKAGE_DST=.

function packaging() {
    echo "Generating ${PACKAGE_NAME} for OS: ${OS} ARCH: ${ARCH}"

    # Building the binary
    GOOS=${OS} GOARCH=${ARCH} make build && \
    # Copying the binary to the current working dir
    cp ${BINARY_DST}${BINARY_NAME} . && \
    # Build the package for the built binary
    tar -czf ${PACKAGE_DST}/${PACKAGE_NAME} ${BINARY_NAME} && \
    # Creating the checksum for the package
    shasum -a 256 -b ${PACKAGE_DST}/${PACKAGE_NAME}

    # Removing the copied binary since it was left on the original location
    rm -f ${BINARY_NAME}
}

# Run packaging procedure
if [ "$1" != "" ] && [ -d $1 ]; then
    PACKAGE_DST=$1
else
    echo -e "The destination $1 is not valid. The default ${PACKAGE_DST} will be taken instead."
fi

packaging
