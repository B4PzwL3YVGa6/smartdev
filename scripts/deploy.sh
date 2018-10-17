#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [ -z "${APC_USERNAME}" ]; then
	echo "APC_USERNAME must be set!"
	exit 1
fi
if [ -z "${APC_PASSWORD}" ]; then
	echo "APC_PASSWORD must be set!"
	exit 1
fi
if [ -z "${APC_ORGANIZATION}" ]; then
	echo "APC_ORGANIZATION must be set!"
	exit 1
fi
if [ -z "${APC_SPACE}" ]; then
	echo "APC_SPACE must be set!"
	exit 1
fi

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
	cd ..
fi
echo $PWD

# =============================================================================================
echo "deploying smartdev app ..."

wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.32.0&source=github-rel' -qO cf-cli.tgz
tar -xvzf cf-cli.tgz 1>/dev/null
chmod +x cf
rm -f cf-cli.tgz || true

./cf login -a "https://api.lyra-836.appcloud.swisscom.com" -u "${APC_USERNAME}" -p "${APC_PASSWORD}" -o "${APC_ORGANIZATION}" -s "${APC_SPACE}"

# make sure routes will be ready
./cf create-route "${APC_SPACE}" scapp.io --hostname smartdev
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname smartdev
./cf create-route "${APC_SPACE}" scapp.io --hostname smartdev-blue-green
./cf create-route "${APC_SPACE}" applicationcloud.io --hostname smartdev-blue-green
sleep 2

# secure working app
./cf rename smartdev smartdev_old || true
./cf unmap-route smartdev_old scapp.io --hostname smartdev-blue-green || true
sleep 2

# push new app
./cf push smartdev_new --no-route
./cf map-route smartdev_new scapp.io --hostname smartdev-blue-green
./cf map-route smartdev_new applicationcloud.io --hostname smartdev-blue-green
sleep 5

# test app
response=$(curl -sIL -w "%{http_code}" -o /dev/null "smartdev-blue-green.scapp.io")
if [[ "${response}" != "200" ]]; then
    ./cf delete -f smartdev_new || true
    echo "App did not respond as expected, HTTP [${response}]"
    exit 1
fi

# finish blue-green deployment of app
./cf delete -f smartdev || true
./cf rename smartdev_new smartdev
./cf map-route smartdev scapp.io --hostname smartdev
./cf map-route smartdev applicationcloud.io --hostname smartde
./cf unmap-route smartdev scapp.io --hostname smartdev-app-blue-green || true
./cf unmap-route smartdev applicationcloud.io --hostname smartdev-blue-green || true
./cf delete -f smartdev_old

# show status
./cf apps
./cf app smartdev

./cf logout

rm -f cf || true
rm -f LICENSE || true
rm -f NOTICE || true
