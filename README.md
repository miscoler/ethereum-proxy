Golang/Rust Test Assignment
Make a caching proxy for eth_getBlockByNumber.

It should expose this api as a set of REST APIs, such as /block/latest/txs/1 or /block/latest/txs/<hash>.

/block/ has either "latest" or a number as a param.

/block/<>/txs/ has either index or tx hash as a param.

The proxy should cache responses in LRU style.

It should use Cloudflare Eth Gateway for that: https://developers.cloudflare.com/distributed-web/ethereum-gateway/interacting-with-the-eth-gateway

Use Go or Rust to implement this assignment.

Keep in mind that the last ~20 blocks in the network could change due to reorgs.

Create a GitHub or GitLab repo for that and share with us.
--------------------------------------------------------------------------------------
Features:
- viper for config
- zap-logger for logging
- prometheus for metrics
- gomock for tests
- golangci passed
- graceful shutdown

Possible ways to improve performance if needed:
- use easyjson;
- use fasthttp;