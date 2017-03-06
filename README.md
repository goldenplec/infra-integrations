# On-host Integrations

## building

To test and build all the integrations:

```bash
$ make test
$ make build
```

Instead, if we want to test and build one of the integrations:
```bash
$ INTEGRATION=nginx make test
$ INTEGRATION=nginx make build
```
