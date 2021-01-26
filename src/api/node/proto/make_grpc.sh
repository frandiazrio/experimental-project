
# sudo ./make_sh to run script



#make sure you have the latest version of protocol buffers and grpc (Pull them if necessary)
eval go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
eval go get -u google.golang.org/grpc
echo Making grpc interface
eval protoc --go_out=plugins=grpc:. \
	--go_opt=paths=source_relative \
	node.proto
echo Done!
#eval protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#      node.proto
