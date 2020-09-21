# go-port-tester
Simple port tester to check connectivity between hosts using specified ports

For example you want to check if tcp/udp port 8000 is reachable between 2 hosts in diffferent networks:

```
parallel-ssh -H 1.1.1.2 -H 2.2.2.2 -i './go-port-tester -proto tcp -port 8000 1.1.1.2 2.2.2.2

[1] 18:50:01 [SUCCESS] 1.1.1.2
1.1.1.2: OK
2.2.2.2: OK
[2] 18:50:02 [FAILURE] 2.2.2.2
2.2.2.2: OK
1.1.1.2: Error connection timeout

parallel-ssh -H 1.1.1.2 -H 2.2.2.2 -i './go-port-tester -proto udp -port 8000 1.1.1.2 2.2.2.2
[1] 18:50:01 [SUCCESS] 1.1.1.2
1.1.1.2: OK
2.2.2.2: OK
[2] 18:50:02 [SUCCESS] 2.2.2.2
2.2.2.2: OK
1.1.1.2: OK
```

parallel-ssh web-page: https://parallel-ssh.org/
