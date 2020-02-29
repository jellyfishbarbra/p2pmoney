# verzo 

A minimalist and experimental decentralized ledger with an identity-based accounting.

### Alpha version

#### Goal

Improve my software engineering and protocol design skills by building a distributed ledger from scratch.

#### Architecture

##### Modules

###### Ledger module

- Define internal data structures to represent the blockchain
- API to store and retrieve ledger state from memory/disk
- Maintain the current state of the ledger (block tip, difficulty, mempool etc.)
- Ignorant of consensus logic

###### Consensus module
- Encapsulates any consensus-breaking logic
- Enforce consensus rules/invariants:
  - Rules for valid block, transaction, data definitions etc.
  - "Business logic" to the ledger module


###### Network module
- Encapsulate all network interactions
- Provides a clean API to the rest of the client
- Peering daemon:
  - Disconver and maintain connections to *compatible* peers
  - Broadcast/gossip data to other peers
  - Maintain ban list of malicious peers
  - Maintain state about the node's view of the network (peers, telemetry)
- JSON-RPC daemon
  - API providing node's state to the outside world
    - Ledger state 
    - Mempool
    - Block information
