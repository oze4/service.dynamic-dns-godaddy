### Homegrown Dynamic DNS for GoDaddy

```yaml
# ...
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
