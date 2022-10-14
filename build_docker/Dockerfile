FROM ubuntu:bionic
RUN sed -i s/archive.ubuntu.com/mirrors.aliyun.com/g /etc/apt/sources.list \
    && apt update \
    && apt install -y ffmpeg
WORKDIR /usr/local/uni-ledger/medical-sys/application/
COPY application/web web/
WORKDIR /usr/local/uni-ledger/medical-sys/
COPY docker/docker-config.yaml ./config.yaml
WORKDIR /data/database
COPY application/database/*.sqlite ./
WORKDIR /usr/bin
COPY build/medical_sys build/medical_sys_tools ./
EXPOSE 80
VOLUME /data
CMD ["medical_sys"]
