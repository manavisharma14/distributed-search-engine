Distributed Search Engine

A fault-tolerant distributed search engine built in Go to explore search infrastructure, distributed systems, and high-performance query processing. The system uses inverted indexing, document sharding, concurrent query execution, gRPC-based shard communication, PostgreSQL persistence, Redis caching, and Dockerized deployment.

Features

* Inverted index for fast full-text search
* Horizontal document sharding
* Concurrent query execution across shards
* Coordinator-shard architecture using gRPC
* PostgreSQL-backed document storage
* Redis query caching
* Fault-tolerant query handling
* Dockerized deployment
* REST API and web interface

Architecture

Client
  |
Coordinator
  |
  +---- Shard 1
  |
  +---- Shard 2
  |
  +---- Shard N

PostgreSQL -> Document Storage
Redis      -> Query Cache


Each shard maintains its own inverted index and processes queries independently. The coordinator fans out requests to all shards, aggregates results, ranks them, and returns the final response.

Performance

* Indexed and searched 250,000+ documents
* Sustained 500+ QPS under 100 concurrent clients
* Reduced p95 query latency from 340ms to 120ms through parallel shard execution
* O(1) document retrieval using in-memory document mapping

Tech Stack

* Go
* gRPC
* PostgreSQL
* Redis
* Docker

Running

git clone https://github.com/<username>/distributed-search-engine.git
cd distributed-search-engine

docker compose up --build

Search API:
GET /search?q=distributed+systems

Key Concepts

* Inverted Indexing
* Sharding
* Distributed Query Execution
* Goroutines, Channels, and WaitGroups
* Fault Tolerance
* Caching
* Coordinator-Worker Architecture
