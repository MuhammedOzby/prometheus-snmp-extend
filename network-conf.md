# Router yapısı

Aşağıdaki network topolojisi için konfigurasyon dosyasıdır. Router kısmı bulunurken SWler için Telnet üzerinden otomatik basılacaktır. LAB için hızlı olması adına. Birde fazla network odaklı değil daha çok test amaçlıdır.

![network topology](doc-files/net-topo.png)

## Linux cihazımızda

IP route için routerlara statik yönlendirme:

> GNS3 altındaki NAT ağı 192.168.122.0/24 şeklindedir ve Gateway: 192.168.122.1 olacaktır.

```bash
sudo ip route add 10.0.0.0/24 via 192.168.122.1 dev virbr0
```
