#!/bin/bash
string=`git describe --abbrev=0`
version=$(echo "$string" | awk -F. '{print $1"."$2"."$3+1}')
echo $version
git add .
git commit -m "update version to $version  $1"
git tag $version -m "$version"
git push --tags
git push
sleep  10s