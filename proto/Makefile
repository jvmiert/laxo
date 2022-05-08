
.PHONY: get gen

get:
	rm -rf ./protoc/
	mkdir ./protoc
	wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win64.zip -O ./protoc/protoc.zip
	tar -xf ./protoc/protoc.zip -C ./protoc/
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	cd ./protoc/ && npm.cmd install ts-protoc-gen@v0.15.0

gen:
	rm -rf ../frontend/proto
	rm -rf ../backend/laxo/proto/gen
	mkdir ../backend/laxo/proto/gen
	mkdir ../frontend/proto
	./protoc/bin/protoc --go_out=../backend/laxo/proto/gen/ --go_opt=paths=source_relative --go-grpc_out=../backend/laxo/proto/gen/ --go-grpc_opt=paths=source_relative ./*.proto
	@PATH="$(PATH):$(CURDIR)/protoc/bin:$(CURDIR)/protoc/node_modules/.bin;" \
	./protoc/bin/protoc --plugin=".\protoc\node_modules\.bin\protoc-gen-ts.cmd" --js_out="import_style=commonjs,binary:." --ts_out="service=grpc-web:." ./*.proto
	mv ./*.js ../frontend/proto
	mv ./*.ts ../frontend/proto
