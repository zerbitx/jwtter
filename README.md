# Jwtter (jotter)

Just a little tool for my own use to make jwts and verify them from the cli.

## new subcommand help
```bash
Uses signing key in config, or environment variable JWT_SIGNING_KEY

Usage:
  jwtter new '{"json":"of", "your": "claims"}' [flags]

Flags:
  -d, --duration duration   -d 120h (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".)
  -h, --help                help for new
  -a, --issuedAt int        -a 1590524782 (default 1590534065)
  -i, --issuer string       -i veridian-dynamics
  -k, --key string          JWT signing key

Global Flags:
      --config string   config file (default is $HOME/.jwtter.yaml)
```


## verify subcommand help
```bash
Uses signing key in config, or environment variable JWT_SIGNING_KEY

Usage:
  jwtter verify [flags]

Flags:
  -h, --help         help for verify
  -k, --key string   JWT signing key

Global Flags:
      --config string   config file (default is $HOME/.jwtter.yaml)

```