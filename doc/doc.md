### Repository Overview

Note: go-resonate is in the early (early) stages of development and is not usable in any way for anything yet.

#### API
The client-side application interface. It's currently designed to run locally and communicate with BigchainDB/IPDB over http, though this may change. Working on functionality for the following user actions:
- Register a user
- Login 
- Create a project
- other TBD

#### Bigchain
[Handcrafting BigchainDB transactions](https://docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html), serializing to json, and sending them to BigchainDB/IPDB.

#### Coala
A json-ld implementation of the Coalaip Rights Reference Model ([RRM](https://github.com/COALAIP/specs/tree/master/data-structure)).

#### Types
Contains the `user` type (below) and others.

```go
type User struct {
	Email    string `json:"email"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	// other TBD
}
```

#### Util

----

### Walkthrough

#### Registering a new user 
What you do:
-  Enter your email, username, password, and account type into the interface

What happens next:
- A `public key` and `private key` are generated for your account 
- Your personal information is signed with your `private key`, producing a `user signature`
- Your `user signature` is sent to BigchainDB/IPDB and stored

What you get:
- Your `public key`, `private key`, and a `user id`, which identifies your `user signature` in BigchainDB/IPDB

Note: none of your personal information is communicated or stored in raw form!

#### Logging in
What you do:
- Enter your `private key`, `user id`, email, username, password, and account type into the interface. Woah, that's a lot of stuff, maybe we can have local file containing some/all of this information that a user uploads and then is parsed by the application?

What happens next:
- The application fetches the `user signature` corresponding to the provided `user id` from BigchainDB/IPDB
- It verifies whether the `user signature` is a valid signature of user information with the provided `public key`
- Your `public key`, `private key`, `user id`, and personal information are kept in memory for the remainder of the session

What you get:
- A friendly welcome message

More to come!