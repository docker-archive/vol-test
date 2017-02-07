#!/bin/bash

echo "Building test environment\n"
echo "Input ssh password:"
read -s password

for n in $(seq 1 3); do
  docker-machine create --driver virtualbox node$n
done;

leader_ip=`docker-machine ip node1`
docker-machine ssh node1 docker searm init --advertise-addr $leader_ip
token=`docker-machine ssh node1 docker swarm join-token worker -q`
for n in $(seq 2 3); do
  docker-machine ssh node$n docker swarm join --token $token $leader_ip:2377
done;
for n in $(seq 1 3); do
  docker-machine ssh node$n docker plugin install --grant-all-permissions vieux/sshfs;
done;
#for n in $(seq 1 3); do
#  docker-machine ssh node$n docker volume create -d vieux/sshfs -o sshcmd=$USER@192.168.99.1:~/tmp -o password=$password testvol;
#done;
