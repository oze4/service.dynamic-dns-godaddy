### Homegrown Dynamic DNS for GoDaddy

Dynamically update all `A` records for a specific domain.

The ENV variable, `BASELINE_RECORD` is the `A` record we use as a "database". This is the record we use to compare with what your actual public IP is. If your actual public IP is different than the value of `BASELINE_RECORD` DNS record, we update all GoDaddy `A` records.

```yaml
# ...removed for brevity
          containers:
            - name: godaddy-dynamic-dns
              image: oze4/godaddy-dynamic-dns:latest
              env:
                - name: GODADDY_APIKEY
                  value: "-"
                - name: GODADDY_APISECRET
                  value: "-"
                - name: GODADDY_DOMAIN
                  value: "-"
                - name: BASELINE_RECORD
                  value: "-"
          restartPolicy: OnFailure
```
