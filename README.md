# weasel
Search wallets with balance positive from randomly private key.

# Use

```sh
docker run mauroalderete/weasel:latest -e THREAD=10 -e GATEWAY=https://cloudflare-eth.com -v $HOME/weasel/store:/app/store
```

# Docker compose

```yaml
version: "3.9"

services:
  weasel:
    image: mauroalderete/weasel:latest
    volumes:
      - $HOME/weasel/store:/app/store 
    environment:
      - THREAD=1
        GATEWAY=https://cloudflare-eth.com
```

## references

[goethereumbook](https://goethereumbook.org/en/)