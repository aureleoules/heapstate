# heapstate
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/powered-by-electricity.svg)](https://forthebadge.com)

> build and deploy your apps effortlessly

## Get started

Deploying `heapstate` itself is very easy.

```
$ git clone https://github.com/aureleoules/heapstate && cd heapstate
$ sudo docker-compose up -d
```

Heapstate will be available at [localhost:6000](http://localhost:6000).
## Development

To run `heapstate` in a development environment:
```
$ sudo docker-compose \ 
    -f docker-compose.yml \
    -f docker-compose.dev.yml \
    up -d
```

Heapstate's api will be listening on port `9000` by default.

## Contributing

This repository contains the code of heapstate's backend. If you wish to contribute to the front-end, go [here](https://github.com/aureleoules/heapstateapp).  

The master branch is regularly built and tested, but is not guaranteed to be completely stable. Tags are created regularly to indicate new official, stable release versions of `heapstate`.

The contribution workflow is described in [CONTRIBUTING.md](CONTRIBUTING.md)

## License

[MIT](LICENSE.md)