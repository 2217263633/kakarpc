# git push origin :refs/tags/vv0.0.102  删除分支
#!/bin/bash
string=`git describe --tags --abbrev=0`
version=$(echo "$string" | awk -F. '{print $1"."$2"."$3+1}')
echo $version
# sleep 10
git add .
git commit -m "update version to $version  $1"
git tag $version -m "$version"
git push --tags
echo "tags finish--------"
git push origin master
sleep  10s