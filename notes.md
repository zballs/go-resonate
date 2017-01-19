Just some notes about what's being built here-- later to be formalized into documention...

### Overview 
Building a music streaming service/cooperative "on the blockchain". Still considering what stack is best for the job.
Here's one lineup we have...
- Coalaip data models for creative works and rights assingment/transfer
- BigchainDB/IPDB for user and metadata storage
- Tendermint for the consensus engine, decentralized application architecture
- ?? for payments 

### Identity
There are two types that are central to identity in go_resonate, `account` and `user`, shown in Go code below. We use the Ed25519 implementation from Tendermint's go-crypto library for key generation, signature, and verification. TODO: make sure Ed25519 system in BigchainDB cryptoconditions library is compatible

```go
type Account struct {
	PubKey *PublicKey
	Sequence int 
}

type User struct {
	Email string 
	Name string
	Type string // artist, listener
	// Other info TBD 
} 
```

`Account` contains a `public_key` used for verification and an incrementing `sequence`, which denotes the number of transactions the account has submitted (i.e. to prevent replay attacks). `Account` bytes (or `account` hashes) are stored by the Tendermint ABCI application in a merkle tree.

`User` contains personal information (email, username, other) and a specified type. When a new user is created, a corresponding `user_signature` is stored in BichainDB. The `user_signature` is generating by signing the `user` bytes with a `private_key`. Therefore, no personal information is stored in the database in raw form.

#### Creating a user/account
Tell IPDB to create a user
- Enter your email, username, password, user type, and other personal information into the front-end interface
- Using bcrypt, a secret is generated from the provided password, and a new keypair is generated from the secret
- A `user` instance is created with information provided, serialized to json bytes, and signed with the private key
- We send a POST request to IPDB with the `user_signature`
- If the `user_signature` is sucessfully persisted to the db, we receive an id for the transaction, which we'll call the `user_id`

Tell Tendermint to create an account
- The `user_signature` is signed with the `private_key` to produce an `account_signature`
- An `action` containing the `account_signature`, `user_id`, and `public_key` is prepared and signed
- The `action` bytes are broadcasted to other nodes in the network
- A node receiving the action asks IPDB if it has a transaction with `user_id`
- If so, the node checks whether the `account_signature` is a valid signature of the `user_signature` from IPDB (the node has the sender's `public_key` to verify this)
- If verification succeeds, the node communicates with the ABCI application to store a new account (or account hash) in the merkle tree

#### Login
TODO

#### Removing a user/account
TODO

### Consensus
#### Proof of Work
Pros 
- Bitcoin is still alive

Cons
- Slower
- Energy inefficient

#### Proof of Stake 
Pros
- Faster
- Maybe better for ordering application actions/transactions as opposed to cryptocurrency 

Cons
- "Nothing at stake" issue 
- How to circumvent FLP impossibility result? Randomize or introduce timeouts. People don't like randomization but honey badger says an adversarial scheduler can thwart PBFT ("When Weak Synchrony Fails").. is this a practical concern?

#### Alternatives - proofs of storage/space(time)
Pros
- Something at stake (i.e. disk space), but depending on the protocol, an adversary may be able to mine a block using significantly less disk space than we want
- Many people have empty disk space on their computers, which they might be willing to dedicate to PoS mining

Cons
- Very few (if any) impls
- Make our own impl, e.g. following specification in Spacemint
- ^Fun, hard, not fun, probably a failure




