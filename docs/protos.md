# Protos: Instructions

This project uses [protocol buffer][protobuf] files to define types shared between the
different services and surfaces of the application. These files must be updated and the output
files regenerated when new fields are added.

**NOTE**: The tool registry used by `buf` has a rate limit of 10 unauthenticated requests
per hour :/. If this quota is reached, use [protoc directly instead][#protoc]. 

## Set up environment

1. Install the [`buf` CLI][buf].

    ```sh
    $ GO111MODULE=on go install github.com/bufbuild/buf/cmd/buf@v1.47.2
    ```

1. Check the installation.

    ```sh
    buf --version
    ```

## Update protobuf files

1. Build the protos. From the `protos/` directory, run the following command.

    ```sh
    buf build
    ```

1. Generate the protos. From the `protos/` directory, run the following command.

    ```sh
    buf generate
    ```

## Generate with protoc-gen-go {#protoc}

1. TODO

[buf]:      https://buf.build/docs/tutorials/getting-started-with-buf-cli
[protobuf]: https://protobuf.dev/getting-started/gotutorial/