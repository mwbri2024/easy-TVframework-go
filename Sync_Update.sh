#!/bin/bash

echo "=========编译二进制文件========="

make linux_amd64
make linux_arm64
make linux_armv7

version=$(grep 'const Version' config/config.go | awk -F'"' '{print $2}')

output_file="build/latest_version.json"
echo "{" > "$output_file"
echo "    \"version\": \"$version\"," >> "$output_file"
echo "    \"hash\": {" >> "$output_file"
echo "=========计算文件MD5值========="

first=true
for file in build/*; do
    if [ -f "$file" ] && ([[ "$(basename "$file")" == "itv_linux_amd64" ]] || [[ "$(basename "$file")" == "itv_linux_arm64" ]] || [[ "$(basename "$file")" == "itv_linux_armv7" ]]); then
        # 计算 MD5
        md5=$(md5sum "$file" | awk '{print $1}')
        
        # 写入 JSON
        echo "        \"$(basename "$file")\": \"${md5}\"," >> "$output_file"

        # 生成 tar.gz 压缩包
        tar -czf "build/$(basename "$file").tar.gz" -C build "$(basename "$file")"
    fi
done

echo "    }" >> "$output_file"
echo "}" >> "$output_file"

# sed -i ':a;N;$!ba;s/,\n}/\n}/' "$output_file"
sed -i ':a;N;$!ba;s/, *\n *}/\n    }/g' "$output_file"

echo "=========推送文件至服务器========="

scp -P 22 build/*.tar.gz build/*.json root@123.123.123.123:/www/wwwroot/app

echo "=========清除临时文件========="

rm build/*.tar.gz build/*.json

echo "=========完毕！========="
