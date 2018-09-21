#!/bin/bash


# This PLATFORMS list is refreshed after every major Go release.
# Though more platforms may be supported (freebsd/386), they have been removed
# from the standard ports/downloads and therefore removed from this list.
#
PLATFORMS="darwin/amd64"
PLATFORMS="$PLATFORMS linux/amd64"


PLATFORMS_ARM="linux freebsd netbsd"

##############################################################
# Shouldn't really need to modify anything below this line.  #
##############################################################


VERSION=`geode version`


type setopt >/dev/null 2>&1


make clean
rm -rf release



SCRIPT_NAME=`basename "$0"`
SOURCE_FILE=`echo $@ | sed 's/\.go//'`
CURRENT_DIRECTORY=${PWD##*/}
WORKDIR="./release" # if no src file given, use current dir name


WORKDIRABS=`realpath $WORKDIR`
GODIRABS=`realpath ./pkg/cmd/geode`



mkdir -p $WORKDIRABS


for PLATFORM in $PLATFORMS; do

  cd $WORKDIRABS

  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}

  NAME="geode-$VERSION-$GOOS-$GOARCH"

  TARGETDIR="$WORKDIRABS/$NAME"
  TARGETBINDIR="$TARGETDIR/usr/local/bin"

  TARNAME="$NAME.tar.gz"

  mkdir -p $TARGETDIR
  mkdir -p $TARGETBINDIR

  BIN_FILENAME="$TARGETBINDIR/geode"

  mkdir -p "$TARGETDIR/usr/local/lib/geodelib"
  cp -a "../lib/" "$TARGETDIR/usr/local/lib/geodelib/"


  cd $GODIRABS
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} $@"
  eval $CMD || exit 1
  cd $WORKDIRABS


  PKGNAME="geode-$VERSION"

  tar -czf $TARNAME -C $NAME .
  rm -rf $TARGETDIR

  if [ $GOOS == "darwin" ]
  then
    echo "Building MacOS pkg distributions"
    fpm -s tar -t osxpkg -p "${PKGNAME}.pkg" $TARNAME
  fi


  if [ $GOOS == "linux" ]
  then
    echo "Building Linux pkg distributions"
    fpm -s tar -t deb -n "geode" -v $VERSION -d "libgc-dev" -d "clang" -p "${PKGNAME}.deb" $TARNAME
  fi

  printf "%.20s %22s\n" "${GOOS}-${GOARCH}" "`realpath $TARNAME`"

done











