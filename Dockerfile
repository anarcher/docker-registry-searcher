FROM scratch
ADD docker-registry-searcher
ENTRYPOINT docker-registry-searcher
