# JWK Dummy

This app is a little webserver that signs claims and serves configuration files so that the JWTs can be validated.

To use -
1. First, run the binary. `make run` will do.  
2. Then, you can post your claims to the server and get them signed. `curl http://localhost:3333/sign -X POST -d '{"msg": "Howdy!"}'` will do.

A couple notes:
1. Each time you restart the server, it will generate a new key to sign tokens with. The original goal of this tool was to run as a test issuer for the duration of a test suite. 
2. Only RSA256 signing keys/algos are currently supported. No JSEs or ECDSA keys will ever be returned. This is fine for issuers - all clients should support this. But it won't work as a good tool for testing clients, only testing applications that want to accept a specific set of claims, and a production issuer cannot be used.
