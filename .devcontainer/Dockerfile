FROM mcr.microsoft.com/devcontainers/go:0-1-bullseye
RUN wget -qO - https://raw.githubusercontent.com/yarnpkg/releases/gh-pages/debian/pubkey.gpg | apt-key add - && \
  apt update && apt install -y yasm
ARG FFMEPG_VERSION=n5.0.2
RUN cd /tmp && \
  wget https://github.com/FFmpeg/FFmpeg/archive/refs/tags/${FFMEPG_VERSION}.tar.gz -O ffmpeg.tar.gz && \
  tar -xvzf ffmpeg.tar.gz && rm ffmpeg.tar.gz && \
  cd FFmpeg-${FFMEPG_VERSION} && \
  ./configure --prefix=/usr/local/ffmpeg --enable-shared && make install &&\
  cd .. && rm -r FFmpeg-${FFMEPG_VERSION}
ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/ffmpeg/lib/pkgconfig/
ENV LD_LIBRARY_PATH=/usr/local/ffmpeg/lib/
ENV PATH=$PATH:/usr/local/ffmpeg/bin/
