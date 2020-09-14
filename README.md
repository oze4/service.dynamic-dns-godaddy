### Homegrown Dynamic DNS for GoDaddy

We use `postgres` - must supply necessary info via `.env` file. See `.env.example` (remove the `.example` after setting values)

If you would like to use the container that is already built (in my repo) you can use a similar `docker-compose` file, or `.yaml` file if you're running in Kubernetes (like I am)..

```yaml
# ...
          containers:
            - name: godaddy-dynamic-dns
              image: oze4/godaddy-dynamic-dns:latest
              env:
                - name: PG_HOST
                  value: "-"
                - name: PG_PORT
                  value: "-"
                - name: PG_USER
                  value: "-"
                - name: PG_PASSWORD
                  value: "-"
                - name: PG_DATABASE
                  value: "-"
                - name: GODADDY_APIKEY
                  value: "-"
                - name: GODADDY_APISECRET
                  value: "-"
                - name: GODADDY_DOMAIN
                  value: "-"
          restartPolicy: OnFailure
```
