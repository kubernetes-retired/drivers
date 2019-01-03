FROM centos:centos7
LABEL maintainers="Kubernetes Authors"
LABEL description="Image Driver"

RUN \
  yum install -y epel-release && \
  yum install -y buildah && \
  yum clean all

COPY ./_output/imageplugin /imageplugin
ENTRYPOINT ["/imageplugin"]
