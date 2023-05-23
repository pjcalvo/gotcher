# Contributing to this repo

## Todo list

- add README + docs
- parameterize authentication based on config
- change config to intercept based on params
- parameterize CORS
- allow polling on the file each 10 seconds instead of a single load

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
- release as a binary (high prio)
- accept an array of configurations instead of one (low prio) (maybe not necessary)
