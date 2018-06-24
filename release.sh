#!/bin/bash


# This PLATFORMS list is refreshed after every major Go release.
# Though more platforms may be supported (freebsd/386), they have been removed
# from the standard ports/downloads and therefore removed from this list.
#
PLATFORMS="darwin/amd64" # amd64 only as of go1.5
PLATFORMS="$PLATFORMS windows/amd64 windows/386"
PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS linux/ppc64 linux/ppc64le"
PLATFORMS="$PLATFORMS linux/mips64 linux/mips64le"
PLATFORMS="$PLATFORMS freebsd/amd64"
PLATFORMS="$PLATFORMS netbsd/amd64"
PLATFORMS="$PLATFORMS openbsd/amd64"
PLATFORMS="$PLATFORMS dragonfly/amd64"
PLATFORMS="$PLATFORMS plan9/amd64 plan9/386"
PLATFORMS="$PLATFORMS solaris/amd64"


PLATFORMS_ARM="linux freebsd netbsd"

##############################################################
# Shouldn't really need to modify anything below this line.  #
##############################################################

type setopt >/dev/null 2>&1

cd "pkg/cmd/geode"

SCRIPT_NAME=`basename "$0"`
SOURCE_FILE=`echo $@ | sed 's/\.go//'`
CURRENT_DIRECTORY=${PWD##*/}
OUTPUT="../../../build" # if no src file given, use current dir name

printf "TARGET           PATH\n"
echo   "========================================"


for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="${OUTPUT}/${GOOS}/${GOARCH}/geode"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} $@"
  printf "%.20s %22s\n" "${GOOS}-${GOARCH}" "build/${GOOS}/${GOARCH}"
  eval $CMD
done
