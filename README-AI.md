# Federasyonlu Prometheus ve Uzman Servislerle L2 Gözlemci LAB Projesi

Bu proje, daha önce teorik olarak tasarlanmış federasyonlu ve Go tabanlı uzman servislere dayalı gözlemci mimarisinin çalışan bir prototipini oluşturmak için tasarlanmış bir LAB projesidir. Temel amacı, ağ cihazlarından (simüle edilmiş L2 switch'ler) SNMP verilerini toplamak, özel Go servisi ile hedef keşfi yapmak ve toplanan metrikleri federasyonlu bir Prometheus kurulumu aracılığıyla Grafana'da görselleştirmektir. Tüm altyapı, hızlı ve izole bir ortam sağlamak için Docker ve Docker Compose kullanılarak simüle edilmiştir.

## Özellikler

*   **Simüle Edilmiş LAB Ortamı:** Docker ve Docker Compose kullanılarak esnek ve kolayca yeniden oluşturulabilir bir ağ izleme LAB ortamı.
*   **Dinamik Keşif Servisi:** Go dilinde yazılmış, ağdaki simüle edilmiş L2 switch'lerin IP adreslerini otomatik olarak keşfeden ve Prometheus için uygun JSON hedef dosyaları (`tüm_switchler.json`) üreten bir servis. Bu sayede manuel hedef yapılandırma ihtiyacı ortadan kalkar.
*   **SNMP Metrik Toplama:** Standart Prometheus `snmp_exporter` kullanılarak L2 switch'lerden SNMP verileri toplanır (örneğin, `ifOperStatus` gibi arayüz operasyonel durumları).
*   **Prometheus Federasyonu:** 5 adet "Leaf" (Yaprak) Prometheus sunucusu (bölgelere göre metrik toplar) ve 1 adet "Global" (Küresel) Prometheus sunucusundan (tüm leaf'lerden özet verileri çeker) oluşan ölçeklenebilir bir izleme mimarisi.
*   **Veri Görselleştirme:** Toplanan tüm SNMP metrikleri, Grafana panoları aracılığıyla zengin görselleştirmelerle sunulur.
*   **Ağ Cihazı Yapılandırma Otomasyonu (Opsiyonel):** GNS3 gibi bir ağ simülasyon aracı üzerinde çalışan Cisco IOS tabanlı anahtarları, Telnet üzerinden otomatik olarak yapılandıran Go tabanlı bir yardımcı araç.

## Teknoloji Yığını

*   **Orkestrasyon:** Docker, Docker Compose
*   **Programlama Dilleri:** Go (Golang)
*   **İzleme ve Gözlemleme:** Prometheus (Federasyon), Grafana, SNMP Exporter
*   **Veri Formatları:** YAML (Yapılandırma dosyaları için), JSON (Hedef keşif dosyaları için)
*   **Ağ Simülasyonu:** GNS3 (Ağ cihazı yapılandırma araçları için, dolaylı), Cisco IOS (Anahtar konfigürasyonları)
*   **Ağ Protokolleri:** SNMP (Veri toplama), Telnet (Ağ cihazı yapılandırma aracı için)
*   **Şablonlama:** Go `html/template` (Ağ cihazı yapılandırma şablonları için)

## Kurulum

Bu proje, ana izleme altyapısını Docker üzerinde çalıştırmaktadır. Ayrıca, GNS3 üzerinde simüle edilen ağ cihazlarını yapılandırmak için ayrı bir Go aracı bulunmaktadır.

### Ön Koşullar

*   [**Docker**](https://docs.docker.com/get-docker/) ve [**Docker Compose**](https://docs.docker.com/compose/install/) yüklü olmalıdır.
*   **Go 1.19+** (GNS3 üzerindeki anahtarları yapılandırmak için `network-confs/edge-configurator` aracını kullanmayı planlıyorsanız gereklidir.)
*   **GNS3** (Ağ cihazlarını simüle etmek ve `network-confs/sw-list.json` dosyasını oluşturmak için gereklidir, bu projenin doğrudan bir parçası değildir ancak LAB kurulumunun bir parçasıdır.)
*   Linux cihazınızda `virbr0` arayüzü üzerinden 192.168.122.0/24 ağına statik bir route eklemeniz gerekebilir (GNS3 NAT ayarları için):
    ```bash
    sudo ip route add 10.0.0.0/24 via 192.168.122.1 dev virbr0
    ```

### Kodu Alma

Proje dosyalarını yerel makinenize indirin:

```bash
git clone <repo_url_buraya> # Klonlamak için uygun URL'yi kullanın
cd l2-observer-lab # Klonlanan dizinin adı bu değilse uygun şekilde değiştirin
```

### Bağımlılıkları Yükleme ve Projeyi Oluşturma

Docker Compose, tüm servis bağımlılıklarını yöneteceğinden, ana proje için ek bir bağımlılık yükleme adımı gerekmez. `docker-compose up -d --build` komutu Go servisini otomatik olarak derleyecektir.

`network-confs/edge-configurator.go` aracını kullanmak isterseniz, bu aracı ayrıca derlemeniz gerekmektedir:

```bash
cd network-confs
go mod tidy # Eğer go.mod dosyası varsa bağımlılıkları indirir
go build -o edge-configurator edge-configurator.go
cd ..
```

## Yapılandırma

Projenin temel yapılandırmaları `docker-compose.yml` dosyası ve ilgili servislerin yapılandırma dizinlerinde (`snmp-exporter/`, `prometheus/`) yer almaktadır. `discovery/targets/tüm_switchler.json` dosyası, `discovery-service` tarafından otomatik olarak oluşturulur ve manuel müdahale gerektirmez.

*   **`snmp-exporter/snmp.yml`:** SNMP metriklerinin nasıl toplanacağını tanımlayan yapılandırma dosyasıdır. Varsayılan olarak temel `if_mib` ayakta kalma durumu kontrolü için yapılandırılmıştır.
*   **`prometheus/leaf-a/prometheus.yml` (ve diğer `leaf-b`, `leaf-c`, `leaf-d`, `leaf-e` Prometheus sunucuları):** Her bir Prometheus leaf örneğinin hedef keşif (file_sd_config kullanarak) ve SNMP Exporter'a yönlendirme (`relabel_configs`) ayarlarını içerir. `prometheus_bolge` etiketi, ilgili leaf'in hangi veri merkezine veya bölgeye ait hedefleri izleyeceğini belirler.
*   **`prometheus/global/prometheus.yml`:** Tüm leaf Prometheus sunucularından `/federate` endpoint'i aracılığıyla özetlenmiş metrikleri çeken ana (global) Prometheus örneğinin yapılandırmasıdır.
*   **`discovery/main.go`:** Simüle edilmiş switch IP adreslerini ve ilgili bölgelerini (örneğin "A_Veri_Merkezi") tanımlayan ve bu bilgileri Prometheus'un okuyacağı `targets/tüm_switchler.json` dosyasına yazan Go kodunu içerir. Bu servis, belirlenen aralıklarla hedef dosyasını günceller.
*   **Ağ Cihazı Konfigürasyonları (`network-confs/` dizini - Opsiyonel):**
    *   `dist-sw-conf.ios` ve `router-conf.ios`: GNS3 ortamındaki dağıtım anahtarları ve yönlendiriciler için Cisco IOS benzeri temel yapılandırma dosyalarıdır.
    *   `edge-sw-conf.ios`: Kenar anahtarları için Go `html/template` ile dinamik IP adresi ve hostname atamaları yapılmasını sağlayan yapılandırma şablonudur.
    *   `sw-list.json`: `edge-configurator` aracı tarafından okunan, GNS3'teki anahtar düğümlerinin konsol portu ve hostname bilgilerini içeren JSON dosyasıdır. Bu dosya genellikle GNS3 projesinden dışa aktarılır.

## Kullanım

### LAB Ortamını Başlatma

Projenin kök dizininde aşağıdaki komutu çalıştırarak tüm Docker konteynerlerini (Keşif Servisi, SNMP Exporter, Prometheus sunucuları ve Grafana) oluşturup başlatabilirsiniz:

```bash
docker-compose up -d --build
```

Bu komut, servisleri arka planda (`-d` parametresi ile daemon modunda) çalıştırır ve gerekli derleme işlemlerini (`--build` parametresi ile, özellikle Go servisi için) yapar.

### Erişim Adresleri

Servisler başlatıldıktan sonra aşağıdaki adreslerden erişebilirsiniz:

*   **Global Prometheus:** `http://localhost:9090`
*   **Grafana:** `http://localhost:3000` (Varsayılan kullanıcı adı/şifre: `admin`/`admin`)
*   **Leaf Prometheus Sunucuları:** `docker ps` komutunu kullanarak veya `docker-compose logs` ile ilgili Prometheus servislerinin maruz bıraktığı portları bulabilirsiniz. Örneğin, `leaf-a` için `http://localhost:XXXX`.

### Doğrulama

1.  **Hedef Dosyası Kontrolü:** `discovery/targets/tüm_switchler.json` dosyasının başarıyla oluşturulduğunu ve güncellendiğini kontrol edin.
2.  **Global Prometheus Hedefleri:** `http://localhost:9090` adresindeki Global Prometheus arayüzünde "Status" -> "Targets" menüsüne gidin. Burada `leaf-a:9090`, `leaf-b:9090` gibi tüm "Leaf Prometheus" sunucularının `federate` işi altında aktif olduğunu görmelisiniz.
3.  **Leaf Prometheus Hedefleri:** Her bir "Leaf Prometheus" sunucusunun (örneğin, tarayıcınızda `leaf-a`'nın portunu bulun) "Status" -> "Targets" menüsünde kendi bölgesine ait simüle edilmiş L2 switch IP adreslerini ve `snmp-exporter:9116` hedefini görmelisiniz.
4.  **Grafana Entegrasyonu:** `http://localhost:3000` adresinden Grafana'ya giriş yapın. Yeni bir Prometheus veri kaynağı ekleyin ve URL olarak `http://global-prometheus:9090` adresini kullanın. `up` gibi basit bir PromQL sorgusu çalıştırarak metriklerin başarıyla geldiğini kontrol edin.

### Ağ Cihazı Yapılandırma Aracını Çalıştırma (Opsiyonel)

Eğer GNS3 ortamında simüle edilmiş Cisco anahtarlarını yapılandırmak isterseniz:

1.  GNS3 projenizin çalıştığından ve yapılandırmak istediğiniz anahtarların Telnet portlarının erişilebilir olduğundan emin olun.
2.  `network-confs` dizinine gidin:
    ```bash
    cd network-confs
    ```
3.  Aracı çalıştırın:
    ```bash
    ./edge-configurator # veya go run edge-configurator.go
    ```
    Bu araç, `sw-list.json` dosyasını okuyacak ve her bir anahtara Telnet üzerinden otomatik olarak `edge-sw-conf.ios` şablonundaki yapılandırmaları uygulayacaktır.

## Servis Olarak Çalıştırma / Dağıtım

Bu proje bir Docker Compose uygulaması olduğundan, servis olarak çalıştırmak için Docker'ın kendi mekanizmalarını kullanmak en uygun ve önerilen yöntemdir.

*   **Daemon Olarak Başlatma:**
    Projenin kök dizininde aşağıdaki komutu kullanarak tüm servisleri arka planda (daemon modunda) başlatabilirsiniz:

    ```bash
    docker-compose up -d
    ```

*   **Servis Durumunu Kontrol Etme:**
    Çalışan Docker konteynerlarının durumunu kontrol etmek için:

    ```bash
    docker-compose ps
    ```
    Tüm servislerin loglarını gerçek zamanlı olarak takip etmek için:
    ```bash
    docker-compose logs -f
    ```

*   **Servisleri Durdurma ve Yeniden Başlatma:**
    Tüm servisleri durdurmak için:
    ```bash
    docker-compose stop
    ```
    Tüm servisleri durdurup tekrar başlatmak için:
    ```bash
    docker-compose restart
    ```
    Tüm servisleri durdurup konteynerları ve ağları kaldırmak için:
    ```bash
    docker-compose down
    ```

*   **Genel Arka Uç Go Servisi Çalıştırma (Eğer Docker Dışında Çalıştırılacaksa):**
    Eğer `discovery-service` gibi bir Go servisini Docker konteyneri dışında, doğrudan bir Linux sunucusunda uzun süre çalışan bir servis olarak çalıştırmanız gerekirse, `systemd` kullanabilirsiniz. Aşağıda temel bir `.service` dosyası şablonu verilmiştir:

    Örneğin, `discovery-service` için `/etc/systemd/system/discovery-service.service` adında bir dosya oluşturun:

    ```ini
    [Unit]
    Description=Go Discovery Service for SNMP Targets
    After=network.target

    [Service]
    ExecStart=/usr/local/bin/discovery-service # Servis uygulamasının tam yolu
    WorkingDirectory=/opt/l2-observer-lab/discovery # Uygulamanın çalışacağı dizin
    Restart=always
    User=your_user # Uygulamayı çalıştıracak kullanıcı adınızı buraya yazın
    Group=your_group # Uygulamayı çalıştıracak grup adınızı buraya yazın
    Environment="PATH=/usr/local/bin:/usr/bin:/bin" # Uygulama için gerekli ortam değişkenleri

    [Install]
    WantedBy=multi-user.target
    ```

    `discovery-service` uygulamasını `/usr/local/bin/` dizinine kopyaladığınızdan ve `WorkingDirectory` yolunu projenizin doğru dizinine (`discovery/` içindeki `targets/` klasörüne yazabilmesi için) ayarladığınızdan emin olun.

    **Yönetim Komutları:**
    *   Servisi başlatma: `sudo systemctl start discovery-service`
    *   Sistem başlangıcında otomatik başlatmayı etkinleştirme: `sudo systemctl enable discovery-service`
    *   Servis durumunu kontrol etme: `sudo systemctl status discovery-service`
    *   Servisi durdurma: `sudo systemctl stop discovery-service`
    *   Servisi yeniden başlatma: `sudo systemctl restart discovery-service`
    *   Servis loglarını görüntüleme: `journalctl -u discovery-service -f`

## API Referansı

Bu proje, kendi özel bir web API'sini harici olarak sunmamaktadır. Prometheus, Grafana ve SNMP Exporter gibi entegre bileşenlerin her birinin kendi API'leri ve web arayüzleri mevcuttur; ancak bu projenin özel kodu (Go keşif servisi veya ağ yapılandırma aracı) bir HTTP API endpoint'i sağlamaz.

## Lisans

Bu proje `LICENSE` dosyasında belirtilen [MIT Lisansı](LICENSE) altında lisanslanmıştır.