#!/usr/bin/env bash
echo "start create trigger:$1"
cd trigger
if [ -d "$1" ]; then
echo "$1 has existed and delete"
rm -rf "$1"
fi

mkdir $1
cd $1
cp -r /home/code/flowgo/core/examples/trigger/* ./
for file in $(pwd)/*; do
    echo $file
    sed -i s/sample/$1/g "$file"
    #sed -i s/sample/$1/g `grep sample -rl --include="$file" ./`
done

export GO111MODULE=on
go mod init github.com/qingcloudhx/flow-plugin/trigger/$1
go mod tidy
echo "finished create trigger:$1"