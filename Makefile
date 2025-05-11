.PHONY: generate
generate:
	protoc -I proto proto/auth.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

.PHONY: docker_build
docker_build:
	docker build -t 'go_jwt_mcs' .

.PHONY: docker_run
docker_run:
	docker run -p 50051:50051 \
	-e DATABASE_URL='postgresql://itmodb_owner:npg_FRgxQfGD3K8w@ep-muddy-violet-a8sqedd0-pooler.eastus2.azure.neon.tech/itmodb?sslmode=require' \
	-e PORT=50051 \
	-e ACCESS_TOKEN_SECRET="akmdl3djnqdoianjdoqn11290201jdi" \
	go_jwt_mcs