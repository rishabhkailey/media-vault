FROM registry.access.redhat.com/ubi9 AS ubi-micro-build
RUN mkdir -p /mnt/rootfs
RUN dnf install --installroot /mnt/rootfs curl --releasever 9 --setopt install_weak_deps=false --nodocs -y && \
    dnf --installroot /mnt/rootfs clean all && \
    rpm --root /mnt/rootfs -e --nodeps setup

FROM quay.io/keycloak/keycloak:23.0.6
COPY --from=ubi-micro-build /mnt/rootfs /
ARG DB_DRIVER=postgres
WORKDIR /opt/keycloak
RUN ./bin/kc.sh build --cache=ispn --cache-stack=kubernetes --db=${DB_DRIVER}
ENTRYPOINT [ "./bin/kc.sh" ]`