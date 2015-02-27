#FROM scratch
FROM debian:wheezy
ADD ./docker-registry-searcher /docker-registry-searcher 
ENTRYPOINT ["/docker-registry-searcher"]
#CMD ["/docker-registry-searcher"]
