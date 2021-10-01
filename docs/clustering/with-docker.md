# Commands


## Create Docker Network

```bash
docker network create tester
```

## Create Volumes

```bash
mkdir /tmp/nats-config
mkdir /tmp/nats-vol-1
mkdir /tmp/nats-vol-2
mkdir /tmp/nats-vol-3
```

## Create Configs for JetStream

```bash
cat << EOF > /tmp/nats-config/cluster-server-1.conf
server_name=n1-c1
listen=14222

jetstream {
   store_dir=/nats/storage
}

cluster {
  name: C1
  listen: localhost:14200
  routes: [
    nats-route://my-jetstream-server-2:24200
    nats-route://my-jetstream-server-3:34200
  ]
}
EOF
```

```bash
cat << EOF > /tmp/nats-config/cluster-server-2.conf
server_name=n2-c1
listen=24222

jetstream {
   store_dir=/nats/storage
}

cluster {
  name: C1
  listen: localhost:24200
  routes: [
    nats-route://my-jetstream-server-1:14200
    nats-route://my-jetstream-server-3:34200
  ]
}
EOF
```

```bash
cat << EOF > /tmp/nats-config/cluster-server-3.conf
server_name=n3-c1
listen=34222

jetstream {
   store_dir=/nats/storage
}

cluster {
  name: C1
  listen: localhost:34200
  routes: [
    nats-route://my-jetstream-server-1:14200
    nats-route://my-jetstream-server-2:24200
  ]
}
EOF
```

## First Docker Conatiner
```bash
docker run \
    -it \
    -p 14222:4222 -p 18222:8222 -p 14200:14200 \
    --rm \
    --name my-jetstream-server-1 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-1,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-1.conf
```


## Second Docker Container

```bash
docker run \
    -it \
    -p 24222:4222 -p 28222:8222 -p 24200:24200 \
    --rm \
    --name my-jetstream-server-2 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-2,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-2.conf
```


## Third Docker Container

```bash
docker run \
    -it \
    -p 34222:4222 -p 38222:8222 -p 34200:34200 \
    --rm \
    --name my-jetstream-server-3 \
    --network tester \
    --mount type=bind,source=/tmp/nats-vol-3,dst=/nats/storage \
    --mount type=bind,source=/tmp/nats-config,dst=/home/nats-config \
    nats:2.6.1 -c /home/nats-config/cluster-server-3.conf
```


