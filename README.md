# launcher

Process supervisor for distroless Java containers. It forks the command in
`LAUNCH_CMD`, forwards SIGTERM/SIGINT to the child, propagates the child's exit
code, and optionally polls the child's gRPC health endpoint so a deadlocked JVM
(alive, not serving) is killed and the container restarted.

Built static and copied into every Java service image as that service's
entrypoint binary (see the `launcher-linux-<arch>` COPY in kernel, eventbus and
stratum Dockerfiles). The fleet's build image carries it at
`/usr/local/bin/launcher`.
