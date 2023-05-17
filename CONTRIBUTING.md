# Contributing to this repo

## Todo list

- add README + docs
- parameterize the port
- parameterize authentication based on config
- change config to patch based on params and METHOD type
- parameterize CORS

```yml
target_url: sometarget
authentication:
    basic: 
        username: something
        password: something
    bearer:
        type: bearer (default to bearer)
        token: token
intercept:
  responses:
    - match:
        uri: "*front/v1/customers"
        params: 
            - name: q
              value: search
        methods: 
            - GET
            - POST
      patch:
        status: 400
        body: ./response.json
        type: file
    - match:
        uri: "*addresses*"
        params: 
            - name: q
              value: search
        methods: 
            - GET
            - POST
      patch:
        status: 400
        type: json
        body: |
            {
            "foo" : "bar"
            }
  requests:
    - match:
        uri: "*front/v1/customers"
        methods: 
            - GET
      patch:
        notfound: true
        reject: false
```

- allow to use single word intercepts (reject: true -> returns a 500, notfound: true -> returns a 404)
- accept verbose as a flag or read it from file (mid prio)
- parameterize the file to be used (easy low prio)
- de-duplicate code from the shoulds functions (high prio)
- release as a binary (high prio)
- accept an array of configurations instead of one (low prio)
