# redpanda
Working with Red Panda

```bash
rpk container start -n 3

Waiting for the cluster to be ready...

Found an existing cluster:

  NODE ID  ADDRESS          
  0        127.0.0.1:54789  
  1        127.0.0.1:54797  
  2        127.0.0.1:54796  
```
Now

```bash
rpk cluster info --brokers 127.0.0.1:54789

rpk topic create cwxstat  --brokers 127.0.0.1:54789
rpk topic list  --brokers 127.0.0.1:54789

rpk topic produce cwxstat  --brokers 127.0.0.1:54789

# Clean up
rpk container purge

```

Ref:
https://github.com/twmb/franz-go
