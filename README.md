# weasel
Search wallets with balance positive from randomly private key.

# Use

```sh
docker run mauroalderete/weasel:latest -e THREAD=10 -v $HOME/weasel.match.json:/app/match.json
```

# Docker compose

```yaml
version: "3.9"

services:
  weasel:
    image: mauroalderete/weasel:latest
    restart: always
    volumes:
      - $HOME/weasel.match.json:/app/match.json
    environments:
      - THREAD=1
```

## references

[goethereumbook](https://goethereumbook.org/en/)