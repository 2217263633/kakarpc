#!/bin/bash
string=`git describe --abbrev=0`
version=$(echo "$string" | awk -F. '{print $1"."$2"."$3+1}')
echo $version
git add .
git commit -m "update version to $version"
git tag $version
git push origin master --tags
sleep  10s