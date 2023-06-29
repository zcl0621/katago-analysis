#!/bin/bash

#use bash -x  /home/zhangchaolong/analysis/katago-server/start.sh katago-1.12.3-v11
# 设置镜像名称和标签
image_name="registry.prd.estargo.com.cn/estartgo-dev/ai/katago-server"
image_tag=$1
container_name_prefix="katago-analysis"

#for (( i=1; i<=72; i++ )); do
#    container_name="${container_name_prefix}-${i}"
#    docker stop "${container_name}" >/dev/null 2>&1
#    docker rm "${container_name}" >/dev/null 2>&1
#done

# 从 32760 开始，为每个容器分配端口
port=32760

# 启动容器
for (( i=1; i<=18; i++ )); do
    container_name="${container_name_prefix}-${i}"
    container_port=$((port + i - 1))
    docker stop "${container_name}" >/dev/null 2>&1
    docker rm "${container_name}" >/dev/null 2>&1
    docker run -d \
        --name "${container_name}" \
	--gpus '"device=0"' \
	--restart always \
	--ulimit nofile=65535:65535 \
        -p "${container_port}:8080" \
        "${image_name}:${image_tag}"
done

sleep 300

for (( i=19; i<=36; i++ )); do
    container_name="${container_name_prefix}-${i}"
    container_port=$((port + i - 1))
    docker stop "${container_name}" >/dev/null 2>&1
    docker rm "${container_name}" >/dev/null 2>&1

    docker run -d \
        --name "${container_name}" \
        --gpus '"device=1"' \
        --restart always \
        --ulimit nofile=65535:65535 \
        -p "${container_port}:8080" \
        "${image_name}:${image_tag}"
done

sleep 300

for (( i=37; i<=54; i++ )); do
    container_name="${container_name_prefix}-${i}"
    container_port=$((port + i - 1))
    docker stop "${container_name}" >/dev/null 2>&1
    docker rm "${container_name}" >/dev/null 2>&1

    docker run -d \
        --name "${container_name}" \
        --gpus '"device=2"' \
        --restart always \
        --ulimit nofile=65535:65535 \
        -p "${container_port}:8080" \
        "${image_name}:${image_tag}"
done

sleep 300

for (( i=55; i<=72; i++ )); do
    container_name="${container_name_prefix}-${i}"
    container_port=$((port + i - 1))
    docker stop "${container_name}" >/dev/null 2>&1
    docker rm "${container_name}" >/dev/null 2>&1

    docker run -d \
        --name "${container_name}" \
        --gpus '"device=3"' \
        --restart always \
        --ulimit nofile=65535:65535 \
        -p "${container_port}:8080" \
        "${image_name}:${image_tag}"
done

echo "容器已启动"
