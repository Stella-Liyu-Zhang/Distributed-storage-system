# Distributed file storage system

Welcome to the comprehensive README for SurfStore, the Distributed file storage system project!

## Introduction

Surfstore is a project that aims to provide a distributed storage system. It utilizes Go and gRPC to implement a server-client architecture for managing files and their metadata across multiple nodes.

## Getting Started

Before diving into the project, it's essential to have a good understanding of the following concepts:

1. Interfaces: Go interfaces are named collections of method signatures. Understanding interfaces is crucial for implementing the Surfstore services. You can refer to the following resources to learn more about interfaces in Go:

   - [Go by Example: Interfaces](https://gobyexample.com/interfaces)
   - [How to Use Interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)

2. gRPC: Surfstore leverages gRPC for communication between clients and servers. It's important to have knowledge of writing gRPC servers and clients in Go. The official [gRPC documentation](https://grpc.io/docs/languages/go/basics/) provides a good resource to get started.

## Protocol Buffers

The Surfstore project utilizes protocol buffers for defining message types and gRPC services. The `SurfStore.proto` file contains the message definitions and service declarations. You can find the details in the file.

To generate the gRPC client and server interfaces from the `.proto` service definition, you need to use the protocol buffer compiler `protoc` with the gRPC Go plugin. The following command generates the required files:

```shell
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/surfstore/SurfStore.proto
```

Running this command generates `SurfStore.pb.go` and `SurfStore_grpc.pb.go` files in the `pkg/surfstore` directory. These files contain the necessary code for working with protocol buffers and gRPC.

## Surfstore Interface

The `SurfstoreInterfaces.go` file contains the interface definitions for the BlockStore and MetadataStore:

- `MetaStoreInterface`: Defines methods for retrieving FileInfoMap, updating file metadata, and getting the BlockStore address.
- `BlockStoreInterface`: Defines methods for getting a block, putting a block, and checking the existence of blocks.

## Implementation

The Surfstore implementation consists of both server and client components.

### Server

The server implementation includes the following files:

- `BlockStore.go`: Provides a skeleton implementation of the `BlockStoreInterface`. You need to implement the methods in this file.
- `MetaStore.go`: Provides a skeleton implementation of the `MetaStoreInterface`. You need to implement the methods in this file.
- `cmd/SurfstoreServerExec/main.go`: Contains the `startServer` method, which you must implement. This method registers the appropriate service (MetaStore, BlockStore, or Both) and starts listening for connections from clients.

### Client

The client implementation includes the following files:

- `SurfstoreRPCClient.go`: Provides the gRPC client stub for the Surfstore server. You need to implement the methods in this file.
- `SurfstoreUtils.go`: Contains the `ClientSync` method, which you need to implement. This method handles the synchronization logic for clients.

## Usage

To run the Surfstore server and client, follow these steps:

1. Run the server using the following command:

   ```shell
   go run cmd/SurfstoreServerExec/main.go -s <service> -p <port> -l -d (BlockStoreAddr*)
   ```

   Replace `<service>` with either `meta`, `block`, or `both` to specify the service provided by the server. `<port>` defines the port number that the server listens on (default is 8080). Use `-l` to configure the server to only listen on localhost, and `-d` to enable log output. `(BlockStoreAddr*)` represents the BlockStore address if `service=both`, and it should be in the format `ip:port`.

2. Run the client using the following command:

   ```shell
   go run cmd/SurfstoreClientExec/main.go -d <meta_addr:port> <base_dir> <block_size>
   ```

   Replace `<meta_addr:port>` with the MetaStore server's address and port. `<base_dir>` is the base directory containing the files you want to sync, and `<block_size>` is the desired block size.

## Examples

Here are some example commands to help you get started:

```shell
go run cmd/SurfstoreServerExec/main.go -s both -p 8081 -l localhost:8081
```

This starts a server that listens only on localhost, port 8081, and services both the BlockStore and MetaStore interfaces.

```shell
# Run the commands below on separate terminals (or nodes)
go run cmd/SurfstoreServerExec/main.go -s block -p 8081 -l
go run cmd/SurfstoreServerExec/main.go -s meta -l localhost:8081
```

The first command starts a server that services only the BlockStore interface and listens only on localhost, port 8081. The second command starts a server that services only the MetaStore interface, listens only on localhost, port 8080, and references the BlockStore as the underlying BlockStore. (Note: If these servers are on separate nodes, use the public IP address and remove `-l`)

From a new terminal or node, run the client using the provided script (after building if using a new node). Create a base directory with files and execute the following command:

```shell
mkdir dataA
cp ~/pic.jpg dataA/
go run cmd/SurfstoreClientExec/main.go server_addr:port dataA 4096
```

This command syncs `pic.jpg` to the server hosted at `server_addr:port`, using `dataA` as the base directory and a block size of 4096 bytes.

From another terminal or node, run the client to sync with the server (after building if using a new node):

```shell
mkdir dataB
go run cmd/SurfstoreClientExec/main.go server_addr:port dataB 4096
ls dataB/
```

You will observe that `pic.jpg` has been synced to this client.

## Makefile

A `Makefile` is provided for your convenience. It offers shortcuts to run the BlockStore and MetaStore servers. Use the following commands:

- Run both BlockStore and MetaStore servers (listening on localhost, port 8081):

  ```shell
  make run-both
  ```

- Run only the BlockStore server (listening on localhost, port 8081):

  ```shell
  make run-blockstore
  ```

- Run only the MetaStore server (listening on localhost, port 8080):
  ```shell
  make run-metastore
  ```
  Feel free to modify and enhance this README to provide more information about your specific implementation and any additional details you think would be helpful for others using
