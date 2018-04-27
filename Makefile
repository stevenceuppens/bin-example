OPENAPI_SPEC_FILE__API = ./openapi/spec/bin-api.yaml
OPENAPI_SPEC_FILE__STORE = ./openapi/spec/bin-store.yaml
OPENAPI_GEN_DIR = ./openapi/gen/
OPENAPI_GEN_DIR__API = ./openapi/gen/bin-api
OPENAPI_GEN_DIR__STORE = ./openapi/gen/bin-store

all: 
	@make gen-openapi
	@make build

clean:
	rm -Rf ./bin
	rm -Rf $(OPENAPI_GEN_DIR)

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/cli ./cmd/bin-cli/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/proxy ./cmd/bin-proxy/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/store ./cmd/bin-store/main.go

gen-openapi:
	mkdir -p $(OPENAPI_GEN_DIR)
	mkdir -p $(OPENAPI_GEN_DIR__API)
	mkdir -p $(OPENAPI_GEN_DIR__STORE)
	swagger generate server -t $(OPENAPI_GEN_DIR__API) -f $(OPENAPI_SPEC_FILE__API) -s server --exclude-main
	swagger generate server -t $(OPENAPI_GEN_DIR__STORE) -f $(OPENAPI_SPEC_FILE__STORE) -s server --exclude-main
	swagger generate client -t $(OPENAPI_GEN_DIR__API) -f $(OPENAPI_SPEC_FILE__API)
	swagger generate client -t $(OPENAPI_GEN_DIR__STORE) -f $(OPENAPI_SPEC_FILE__STORE)