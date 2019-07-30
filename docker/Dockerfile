FROM docker.io/centos:7
MAINTAINER beekuaibao

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN yum install -y epel-release
RUN rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
RUN rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm
RUN yum install -y ffmpeg

#RUN yum install -y java
#RUN yum install -y xorg-x11-server-Xvfb chromium

COPY config /app/config
COPY images /app/images
#COPY tools /app/tools
#COPY tools/linux64 /app/tools/linux64
#RUN chmod +x /app/tools/linux64/chromedriver
COPY video-spider /app/
WORKDIR /app

ENTRYPOINT ["/app/video-spider"]