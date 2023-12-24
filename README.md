# common-lib

common-lib is a Golang library designed for the Arcology Network, providing a range of tools and utilities used across different modules within the system.

## Packages 📦

### common 🖇

The `common` package comprises utility and generic functions used in different project modules.

- **generic.go**: Generic functions applicable to different types of data  🔥

- **Parallel Utilities**: Tools for parallel processing, efficiently distributing work among multiple ranges and assigning each range to a dedicated worker function.

- **Basic Math Functions**: Some mathematical operations. 

### codec 🖨

An efficient encoding/decoding library focused on performance and parallelism, primarily used for internal inter-process communication among modules.

### addrcompression 💿
Fast data compression is achieved using a lookup table, replacing addresses with corresponding index numbers.

### mempool 📜

Ensuring thread safety, Mempool is responsible for managing a pool of objects of the same type.

### container 🗄️

This package introduces custom data structures tailored for memory/storage optimization and concurrent uses.

- **pagedArray**: A specialized data structure representing an array divided into multiple blocks or pages.

- **concurrentMap**: Implementation of a concurrent map allowing multiple goroutines to access and modify the map concurrently.

- **orderedSet**: A collection preserving the order of insertion.

### cachedstorage 

This package facilitates Files DBs, memory DB, and offers wrappers for various third-party DB implementations.

- **datastore**: A high-level datastore designed to handle state transition persistency.

- **filedb**: A file database.

- **memDB**: A high-performance DB for concurrent use utilizing a concurrent map.

- **bager DB**: BagerDB wrapper.


## Usage

Include detailed instructions and examples to assist users in integrating and utilizing the common-lib library within their projects.


## License

This project is licensed under the [MIT License](LICENSE).
