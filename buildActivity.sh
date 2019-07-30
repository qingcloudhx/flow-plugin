#!/usr/bin/env bash
#address="github.com/qingcloudhx/flow-plugin/activity/"
#act=${address}$1
#echo ${act}
echo "start create activity:$1"
cd activity
if [ -d "$1" ]; then
echo "$1 has existed and delete"
rm -rf "$1"
fi

mkdir $1
cd $1
cp -r /home/code/flowgo/core/examples/activity/* ./
for file in $(pwd)/*; do
    echo $file
    sed -i s/sample/$1/g "$file"
    #sed -i s/sample/$1/g `grep sample -rl --include="$file" ./`
done

export GO111MODULE=on
go mod init github.com/qingcloudhx/flow-plugin/activity/$1
go mod tidy
echo "finished create activity:$1"