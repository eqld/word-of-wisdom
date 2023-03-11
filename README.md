# Word of Wisdom server and client

Word of Wisdom TCP server and client with PoW verification.

## Server

Server listens for incoming TCP connections on port `8080` and serves a random quote to clients that pass a proof-of-work (PoW) challenge.

The program works as follows:

* The program listens for incoming connections on port 8080 using `net.Listen`.
* For each incoming connection, the program generates a PoW challenge consisting of a random message and a difficulty level. The challenge is sent to the client as a string in the format `<message>:<difficulty>\n`.
* The client must solve the PoW challenge by appending 8 bytes to the message and computing the SHA256 hash of the resulting string. If the first `difficulty` bytes of the hash are all zero, the client's solution is considered valid.
* If the client's solution is valid, the program generates a random quote using the `fortune` command and sends it to the client as a string.

### Usage

TODO

## Client

Command-line client connects to a Word of Wisdom TCP server and receives a random quote from it after successfully solving a Proof of Work challenge.

The program works as follows:

* It first checks if the user has provided the host and port number of the server as command-line arguments. If not, it logs an error message and exits.
* Then it establishes a TCP connection with the server and receives a challenge with difficulty level from the server.
* It parses the challenge with the difficulty level and starts solving the challenge using the Proof of Work algorithm. It generates random strings until it finds a string whose SHA256 hash has a prefix of `difficulty` number of zeros.
* Once it finds the solution, it sends it to the server and receives a quote from the server, which is printed to the STDOUT.

### Usage

TODO

## PoW algorithm selection explanation

Given PoW algorithm was choosen for several reasons:

* Simplicity: this algorithm is relatively simple and easy to understand. It only requires basic operations like concatenation, hashing, and bit manipulation. This makes it easy to implement and verify.

* Security: this algorithm is based on cryptographic hash functions, which are designed to be secure against various attacks. By requiring a certain number of leading zero bits in the hash, it is harder for attackers to generate valid values.

* Resource efficiency: this PoW algorithm is relatively lightweight and doesn't require a lot of computational resources to verify. This makes it suitable for use in a TCP server where it is needed to process many requests quickly.

* Widely used: this PoW algorithm is similar to the one used in Bitcoin and many other cryptocurrencies.

Overall, this PoW algorithm provides a good balance between security and efficiency. While it's not perfect and can still be vulnerable to certain attacks, it's a good starting point for protecting a TCP server from DDoS attacks.
