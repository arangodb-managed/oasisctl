#!/bin/bash

SRCDIR=$(pwd)
REPODIR=${SRCDIR}/.oasis-tools

if [ -z "${GITHUB_USERNAME}" ]; then
    echo GITHUB_USERNAME is not set
    exit 1
fi
if [ -z "${GITHUB_TOKEN}" ]; then
    echo GITHUB_TOKEN is not set
    exit 1
fi
if [ -z "${GITHUB_EMAIL}" ]; then
    echo GITHUB_EMAIL is not set
    exit 1
fi
if [ -z "${VERSION}" ]; then
    echo VERSION is not set
    exit 1
fi
if [ -z "${COMMIT}" ]; then
    echo COMMIT is not set
    exit 1
fi

mkdir -p ${REPODIR}
cd ${REPODIR}
rm -Rf oasis-tools
git clone "https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@github.com/arangodb/oasis-tools.git"
cd oasis-tools

git config user.email "${GITHUB_EMAIL}"
git config user.name "${GITHUB_USERNAME}"

TARGETDIR=oasisctl/${VERSION}-${COMMIT}
mkdir -p ${TARGETDIR}
cp -r -f ${SRCDIR}/bin/* ${TARGETDIR}
ln -sf ${VERSION}-${COMMIT} oasisctl/latest

git add .
git commit -m "Updating oasisctl version ${VERSION}, build ${COMMIT}"
git push

cd ${SRCDIR}