# Commands

## Create Temporary Directory

```bash
mkdir /tmp/nats-js-cluster-setup
mkdir /tmp/nats-js-server-1
mkdir /tmp/nats-js-server-2
mkdir /tmp/nats-js-server-3
```

## KinD Configuration

```bash
cat << EOF > /tmp/nats-js-cluster-setup/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  # https://github.com/kubernetes-sigs/kind/releases
  - role: control-plane
    image: kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
  - role: worker
    image: kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraMounts:
      - hostPath: /tmp/nats-js-server-1
        containerPath: /data/nats-js
  - role: worker
    image: kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraMounts:
      - hostPath: /tmp/nats-js-server-2
        containerPath: /data/nats-js
  - role: worker
    image: kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
    extraMounts:
      - hostPath: /tmp/nats-js-server-3
        containerPath: /data/nats-js

kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        authorization-mode: "AlwaysAllow"
EOF
```

## Start KinD Cluster

```bash
kind create cluster --config /tmp/nats-js-cluster-setup/kind-config.yaml --name nats-js-cluster-test
```

