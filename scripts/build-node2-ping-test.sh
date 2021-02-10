echo "Building Node2 for Ping Test:"
# Move to project root
mv dockerfiles/Make_Node2_Ping_Dockerfile ../Dockerfile

docker build -t node2-ping-test ../
#Put dockerfile back in its place and with its original name
mv ../Dockerfile dockerfiles/Make_Node2_Ping_Dockerfile

echo Cleaning up dangling images

sudo docker rmi -f $(docker images -f "dangling=true" -q)
