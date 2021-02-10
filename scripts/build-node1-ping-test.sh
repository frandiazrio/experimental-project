echo "Building Node1 for Ping Test:"
# Move to project root
mv dockerfiles/Make_Node1_Ping_Dockerfile ../Dockerfile

docker build -t node1-ping-test ../
#Put dockerfile back in its place and with its original name
mv ../Dockerfile dockerfiles/Make_Node1_Ping_Dockerfile

docker rmi -f $(docker images -f "dangling=true" -q)

