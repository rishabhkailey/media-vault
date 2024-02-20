# https://github.com/willhallonline/docker-ansible/blob/master/ansible-core/debian-bullseye-slim/Dockerfile
FROM debian:bookworm

ARG ANSIBLE_CORE_VERSION=2.16.3
ARG ANSIBLE_VERSION=9.2.0
ARG ANSIBLE_LINT=24.2.0
ENV ANSIBLE_CORE_VERSION ${ANSIBLE_CORE_VERSION}
ENV ANSIBLE_VERSION ${ANSIBLE_VERSION}
ENV ANSIBLE_LINT ${ANSIBLE_LINT}

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
  apt-get install -y python3-pip sshpass git openssh-client libhdf5-dev libssl-dev libffi-dev && \
  rm -rf /var/lib/apt/lists/* && \
  apt-get clean

RUN pip3 install --break-system-packages --upgrade pip cffi && \
  pip3 install --break-system-packages ansible-core==${ANSIBLE_CORE_VERSION} && \
  pip3 install --break-system-packages ansible==${ANSIBLE_VERSION} ansible-lint==${ANSIBLE_LINT} && \
  pip3 install --break-system-packages mitogen jmespath && \
  pip install --break-system-packages --upgrade pywinrm && \
  rm -rf /root/.cache/pip

RUN mkdir /ansible && \
  mkdir -p /etc/ansible && \
  echo 'localhost' > /etc/ansible/hosts

WORKDIR /ansible

CMD [ "ansible-playbook", "--version" ]

# FROM python:3.13.0a4-bullseye
# RUN apt update && apt install -y build-essential libssl-dev libffi-dev && \
#   python3 -m pip install --user ansible pyopenssl cryptography
# ENV PATH=$PATH:/root/.local/bin/