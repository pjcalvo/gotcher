# rigo

(Request interceptor in GO) In escense `rigo` is a reverse *Proxy* cli with the ability to intercept and modify requests and responses defined on a configuration file.

## Use cases

Rigo becomes really useful when testing apps or services that are interconnected or dependent on external services. For example:

- Mock an entire backend with custom responses. This would allow testing services with limited connectivity or ahead of development.
- Intercept specific responses that change the behavior of the system under test. This allows to modify and test the behavior of the system under test withouht changing the behavior of it's dependencies.

## Install rigo

```bash
# work in progress
```

### Running rigo

Running rigo is as simple as execute the cli app passing the flags:

```bash
./rigo -p 8400 -f rigo.yaml
```

### Flags

| Flag | Name  | Description | Default |
|------|-------|-------------| ------------|
| p    | port  | Port number to run the proxy on | 8400
| f    | file  | Configuration file to use | rigo.yml
| v    | verbose  | To make the logs more verbose (useful when dealing with CORS) | false
| record    | record  | Record request/response instead of patching | false

### Configuration file

This is the configuration file
