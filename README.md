# Simple Setup of NATS JetStream

**🐳 IMPORTANT NOTE 🐳**

This is a personal note of how I understand NATS JetStream offerings.

Most of the information is based on my understanding, and if any wording or information does not make sense or match the official documentation, please raise an issue or PR.

## Documents

- [Docker based setup](https://github.com/rytswd/simple-nats-js/tree/main/docs/docker-based/README.md) - This only requires Docker, and you can see NATS JetStream in action. You can use CLI to interact with the server.
- [KinD + Helm based setup](https://github.com/rytswd/simple-nats-js/tree/main/docs/kind-based/README.md) - This uses a custom Helm Chart to create a NATS JetStream Cluster. As a part of getting started, this also sets up local Kubernetes cluster with KinD to ease the testing.
- [NACK based setup](https://github.com/rytswd/simple-nats-js/tree/main/docs/nack-based/README.md) - This installs NATS Server and NATS JetStream controller, so that you have a fully running Kubernetes cluster. This also sets up tools with which you can interact with NATS Servers for JetStream features.
