FROM centos:latest 

RUN yum -y install nfs-utils && yum -y install epel-release && yum -y install jq && yum clean all

# Copy nfsplugin from build _output directory
COPY ./_output/nfsplugin /nfsplugin

ENTRYPOINT ["/nfsplugin"]
