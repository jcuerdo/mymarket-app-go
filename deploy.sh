#!/usr/bin/env bash

#WEB

echo 'Deploying WEB'

echo 'Deleting current executable'
rm main

echo 'Downloading new version'
wget 'https://github.com/jcuerdo/mymarket-app-go/releases/download/0.0.1/main'

echo 'Setting privileges'
chmod 755 main

echo 'Killing current'
kill -9 $(lsof -t -i:8080)

echo 'Start service'
nohup ./main < /dev/null > std.out 2> std.err &

#JOBS

echo 'Deploying JOBS'

echo 'Deleting current executable'
rm jobs

echo 'Downloading new version'
wget 'https://github.com/jcuerdo/mymarket-app-go/releases/download/0.0.1/jobs'

echo 'Setting privileges'

chmod 755 jobs

