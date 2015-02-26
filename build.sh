go build
docker build -t anarcher/docker-registry-searcher:0.1.0 .
docker push anarcher/docker-registry-searcher:0.1.0
rm docker-registry-searcher
