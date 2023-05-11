#/bin/bash
pushd $(dirname $0)

make container
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
./out/$GOOS/$GOARCH/vmware-desktop-autoscaler \
    --request-timeout=120s \
    --config=${HOME}/Projects/autoscaled-masterkube-vmware/config/desktop-ca-k8s/config/kubernetes-vmware-autoscaler.json \
    --save=${HOME}/Projects/autoscaled-masterkube-vmware/config/desktop-ca-k8s/config/autoscaler-state.json \
    --log-level=info
