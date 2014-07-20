FROM debian:jessie
ADD ./certs /certs
ADD bastion /usr/local/bin/bastion
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/bastion"]
CMD []
