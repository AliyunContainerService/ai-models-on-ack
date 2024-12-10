## How to generate CRD and clients

cd ./apis && go mod tidy && go mod vendor && bash ./hack/update-crd.sh && bash ./hack/update-codegen.sh