ARG BASE_OS_VERSION=5.13.0-14
ARG BASE_OS_URL=armdocker.rnd.ericsson.se/proj-ldc/common_base_os_release/sles
ARG MODULE_NAME

FROM ${BASE_OS_URL}:${BASE_OS_VERSION}

ARG BASE_OS_VERSION
ARG BASE_OS_REPO=arm.sero.gic.ericsson.se/artifactory/proj-ldc-repo-rpm-local/common_base_os/sles

RUN zypper addrepo --gpgcheck-strict -f https://${BASE_OS_REPO}/${BASE_OS_VERSION} baseos \
        && zypper --gpg-auto-import-keys refresh \
        && zypper install -l -y curl \
	&& zypper install -l -y nodejs \
        && zypper install -l -y npm \
        && npm install -g @apidevtools/swagger-cli \
	&& zypper install -l -y wget \
        && zypper clean --all

COPY app/communication /
CMD ["/communication"]