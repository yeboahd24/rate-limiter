# Algorithms for rate limiting

## Token bucket algorithm

The token bucket algorithm is widely used for rate limiting. It is simple, well understood and
commonly used by internet companies. Both Amazon and Stripe use this algorithm to
throttle their API requests.

The token bucket algorithm work as follows:

- A token bucket is a container that has pre-defined capacity. Tokens are put in the bucket
  at preset rates periodically. Once the bucket is full, no more tokens are added.
- Each request consumes one token. When a request arrives, we check if there are enough
  tokens in the bucket
- If there are enough tokens, we take one token out for each request, and the request
  goes through.
- If there are not enough tokens, the request is dropped.

The token bucket algorithm takes two parameters:

- Bucket size: the maximum number of tokens allowed in the bucket
- Refill rate: number of tokens put into the bucket every second
