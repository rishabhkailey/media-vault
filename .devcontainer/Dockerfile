FROM mcr.microsoft.com/vscode/devcontainers/go:1-bullseye
# RUN wget -qO - https://raw.githubusercontent.com/yarnpkg/releases/gh-pages/debian/pubkey.gpg | apt-key add - && \
#   apt update && apt install -y yasm libvorbis-dev=1.3.7-1 libvpx-dev=1.9.0-1 libx264-dev=2:0.160.3011+gitcde9a93-2.1 libx265-dev=3.4-2 libopus-dev=1.3.1-0.1
# ARG FFMEPG_VERSION=n5.0.2
# RUN cd /tmp && \
#   wget https://github.com/FFmpeg/FFmpeg/archive/refs/tags/${FFMEPG_VERSION}.tar.gz -O ffmpeg.tar.gz && \
#   tar -xvzf ffmpeg.tar.gz && rm ffmpeg.tar.gz && \
#   cd FFmpeg-${FFMEPG_VERSION} && \
#   ./configure --prefix=/usr/local/ffmpeg --enable-shared --enable-libvpx --enable-libvorbis --enable-gpl --enable-libx265 --enable-libx264 --enable-libopus && make install &&\
#   cd .. && rm -r FFmpeg-${FFMEPG_VERSION}
# ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/ffmpeg/lib/pkgconfig/
# ENV LD_LIBRARY_PATH=/usr/local/ffmpeg/lib/
# ENV PATH=$PATH:/usr/local/ffmpeg/bin/
