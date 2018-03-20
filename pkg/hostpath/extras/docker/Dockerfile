FROM alpine
LABEL maintainers="Kubernetes Authors"
LABEL description="HostPath CSI Plugin"

COPY ./_output/hostpathplugin /hostpathplugin
ENTRYPOINT ["/hostpathplugin"]
