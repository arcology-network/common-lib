<h1> common-lib <img align="center" height="50" src="./img/code-circle.svg">  </h1>

common-lib is a Golang library designed for the Arcology Network, providing a range of tools and utilities used across different modules within the system.


<h2> common <img align="center" height="32" src="./img/creative-commons.svg">  </h2>

The `common` package comprises utility and generic functions used in different project modules.

- **generic.go**: Generic functions applicable to different types of data  🔥

- **Parallel Utilities**: Tools for parallel processing, efficiently distributing work among multiple ranges and assigning each range to a dedicated worker function.

- **Basic Math Functions**: Some mathematical operations. 🧮


<h2> codec <img align="center" height="32" src="./img/circle-top-up.svg">  </h2>

An efficient encoding/decoding library focused on performance and parallelism, primarily used for internal inter-process communication among modules.

<h2> addrcompression  <img align="center" height="32" src="./img/archive-down.svg">  </h2>

Fast data compression is achieved using a lookup table, replacing addresses with corresponding index numbers.

<h2> mempool  <img align="center" height="32" src="./img/copy.svg">  </h2>

Ensuring thread safety, Mempool is responsible for managing a pool of objects of the same type.

<h2> container  <img align="center" height="32" src="./img/layers-minimalistic.svg">  </h2>

This package introduces custom data structures tailored for memory/storage optimization and concurrent uses.

| Package           | Description                                                               |
|-------------------|---------------------------------------------------------------------------|
| datastore         | A high-level datastore designed to handle state transition persistency.   |
| filedb            | A file database.                                                          |
| memDB             | A high-performance DB for concurrent use utilizing a concurrent map.      |
| bager DB          | BagerDB wrapper.  


<h2> storage  <img align="center" height="32" src="./img/database.svg">  </h2>

This package facilitates Files DBs, memory DB, and offers wrappers for various third-party DB implementations.

| Package     | Description                                                 |
|-------------|-------------------------------------------------------------|
| datastore   | A high-level datastore designed to handle state transition persistency. |
| filedb      | A file database.                                            |
| memDB       | A high-performance DB for concurrent use utilizing a concurrent map. |
| bager DB    | A [BagerDB](https://github.com/dgraph-io/badger) wrapper.                                           |


<h2> Usage  <img align="center" height="32" src="./img/ruler-cross-pen.svg">  </h2>

Include detailed instructions and examples to assist users in integrating and utilizing the common-lib library within their projects.


<h2> License  <img align="center" height="32" src="./img/copyright.svg">  </h2>

This project is licensed under the [MIT License](LICENSE).
