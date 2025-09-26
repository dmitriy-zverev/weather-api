docker stop $(docker ps -aq) && docker rm $(docker ps -aq)

./scripts/run_docker.sh