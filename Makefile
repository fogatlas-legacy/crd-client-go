packages:=$(shell go list ./pkg/apis... | grep -v /vendor/)

build:
	@ ./hack/update-codegen.sh
	@ ./hack/gen-crd.sh
	@ go build

unit:
	@ go fmt $(packages)
	@ go vet $(packages)
	@ golint $(packages)
	@ shadow $(packages)
	@ go test -cover ./pkg/apis/fogatlas/v1alpha1
