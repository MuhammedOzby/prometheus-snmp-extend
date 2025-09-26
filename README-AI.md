Bu proje kodu ve ek bilgiler doğrultusunda, kapsamlı bir `README.md` dosyası aşağıda Türkçe olarak oluşturulmuştur:

---

# Prometheus SNMP Exporter ile Federasyonlu Ağ Gözlemleme LAB Ortamı

Bu proje, teorik olarak tasarlanmış federasyonlu ve uzman Go servislerine dayalı gözlemci mimarisinin çalışan bir laboratuvar ortamı prototipini oluşturmayı amaçlamaktadır. Ağ cihazlarından SNMP verilerini toplamak, bu verileri merkezi olarak birleştirmek ve görselleştirmek için Go ile geliştirilmiş bir keşif servisi, Prometheus federasyonu, Prometheus SNMP Exporter ve Grafana gibi modern gözlemleme araçlarını bir araya getirir.

## Özellikler

*   **Ağ Cihazı Simülasyonu:** Docker konteynerleri üzerinde çalışan `snmpd` ile temel L2 switch'ler simüle edilir.
*   **Dinamik Keşif Servisi (Go):** Ağdaki hedef IP adreslerini (simüle edilmiş switch'ler) otomatik olarak keşfeder ve Prometheus tarafından okunabilir bir JSON dosyasına (`tüm_switchler.json`) kaydeder. Bu servis, belirli aralıklarla çalışacak şekilde yapılandırılmıştır.
*   **Prometheus SNMP Exporter Entegrasyonu:** Standart Prometheus `snmp_exporter` kullanılarak ağ cihazlarından SNMP metrikleri toplanır.
*   **Prometheus Federasyonu:**
    *   5 adet "Leaf" (Yaprak) Prometheus sunucusu, bölgesel bazda veri toplama ve ön işleme yapar.
    *   1 adet "Global" Prometheus sunucusu, tüm Leaf Prometheus'lardan özet verileri çeker ve merkezi bir görünüm sunar.
*   **Grafana ile Görselleştirme:** Toplanan tüm metrikler, kapsamlı panolar ve uyarılar için Grafana üzerinde görselleştirilebilir.
*   **Docker ve Docker Compose ile Yönetim:** Tüm LAB ortamı bileşenleri (Prometheus instanceleri, SNMP Exporter, Keşif Servisi, Grafana, simüle edilmiş switch'ler), Docker Compose kullanılarak kolayca dağıtılır ve yönetilir.
*   **Ağ Cihazı Konfigürasyon Aracı (Opsiyonel):** GNS3 ortamında gerçek veya sanal Cisco ağ cihazlarının (DIST-SW'ler ve SW-X0/X1/X2/X3/X4 gibi Edge switch'ler) temel IP ve VLAN/EtherChannel konfigürasyonunu otomatik olarak yapmak için Go ile yazılmış bir yardımcı araç.

## Teknoloji Yığını

*   **Programlama Dilleri:** Go (Keşif Servisi ve Ağ Konfigürasyon Aracı için)
*   **Konteynerleştirme:** Docker, Docker Compose
*   **Gözlemleme:** Prometheus (Sunucular ve SNMP Exporter), Grafana
*   **Ağ Simülasyonu:** `snmpd` (Docker imajı)
*   **Ağ Cihazı Yönetimi:** Cisco IOS (konfigürasyon şablonları), Telnet (Go aracı ile otomasyon)
*   **Veri Formatları:** JSON, YAML
*   **Diagram Araçları:** Mermaid (geliştirme sürecinde mimari görselleştirme için)

## Kurulum

### Ön Koşullar

Bu projeyi yerel makinenizde çalıştırmak için aşağıdaki yazılımların kurulu olması gerekmektedir:

*   **Docker Engine:** `v20.10+` veya üzeri önerilir.
*   **Docker Compose:** `v2.0+` veya üzeri önerilir.
*   **Go:** `v1.19+` veya üzeri önerilir (Ağ Konfigürasyon Aracı için).
*   **Git:** Proje kodunu klonlamak için.
*   **(Opsiyonel) GNS3:** Gerçek veya sanal Cisco anahtarlarını yapılandırmak istiyorsanız.

### Kodu İndirme

Proje kodunu yerel makinenize indirmek için aşağıdaki Git komutunu kullanın:

```bash
git clone <repository_url>
cd l2-observer-lab
```
> `<repository_url>` yerine projenin Git deposu adresini yazın.

### Bağımlılıkları Yükleme

*   **Ana LAB Ortamı (Docker Compose ile):** Tüm gerekli Docker imajları (`prom/snmp-exporter`, `prom/prometheus`, `grafana/grafana`, `roenvan/snmpd`) `docker-compose up -d --build` komutu ile otomatik olarak indirilecek veya Go Keşif Servisi için derlenecektir. Manuel bir bağımlılık yüklemesi gerekmez.
*   **Edge Konfigürasyon Aracı (Opsiyonel):** Eğer `network-confs/edge-configurator.go` aracını kullanacaksanız, `network-confs` dizinine gidin ve Go modüllerini indirin:
    ```bash
    cd network-confs
    go mod tidy
    ```

## Yapılandırma

Bu proje, ana `README.md` dosyasında açıklandığı gibi belirli yapılandırma dosyaları içerir.

*   **Keşif Servisi (`discovery/main.go`):**
    *   Simüle edilen switch'lerin IP adresleri (`172.20.0.10`, `172.20.0.11`, `172.20.0.20`, `172.20.0.30`) ve bunlara atanan bölgeler (`A_Veri_Merkezi`, `B_Veri_Merkezi`, `C_Veri_Merkezi`) `main.go` dosyası içinde tanımlanmıştır. Bu değerleri ihtiyaca göre güncelleyebilirsiniz.
    *   Servis, keşfettiği hedefleri `discovery/targets/tüm_switchler.json` dosyasına otomatik olarak yazacaktır. Bu dosya, Prometheus Leaf sunucuları tarafından hedef listesi olarak okunur.

*   **SNMP Exporter (`snmp-exporter/snmp.yml`):**
    *   Bu dosya, SNMP Exporter'ın hangi MIB'leri kullanarak metrik toplayacağını tanımlar. Varsayılan olarak `if_mib` (arayüz durumları) konfigüre edilmiştir. Özelleştirmek için bu dosyayı düzenleyebilirsiniz.

*   **Prometheus Yapılandırmaları (`prometheus/` dizini):**
    *   **Leaf Prometheus (`prometheus/leaf-a/prometheus.yml` vb.):** Her bir Leaf Prometheus yapılandırması, `prometheus_bolge` etiketi ile kendi bölgesini belirtir ve `relabel_configs` kullanarak sadece bu bölgeye ait hedefleri (`tüm_switchler.json` dosyasından) toplar. Ayrıca, hedefleri SNMP Exporter'a (`snmp-exporter:9116`) yönlendirme kurallarını içerir. Diğer Leaf'ler için `prometheus_bolge` ve `regex` değerlerini kendi bölgelerine göre kopyalayıp düzenlemeniz gerekmektedir (örneğin, `leaf-b` için `B_Veri_Merkezi`).
    *   **Global Prometheus (`prometheus/global/prometheus.yml`):** Bu yapılandırma, tüm Leaf Prometheus sunucularından `/federate` endpoint'i aracılığıyla metrikleri çeker. `match[]` parametresi ile hangi metriklerin çekileceği belirtilmiştir (`{job="snmp-switchler"}`).

*   **Docker Compose (`docker-compose.yml`):**
    *   Tüm servislerin (discovery-service, snmp-exporter, simüle edilmiş switch'ler, leaf-a/b/c/d/e, global-prometheus, grafana) tanımlarını, ağ ayarlarını, volume bağlamalarını ve port eşleştirmelerini içerir.
    *   `lab_net` adlı özel bir bridge ağı tanımlanmıştır (`172.20.0.0/16` alt ağı ile).

*   **Ağ Cihazı Konfigürasyon Şablonları (`network-confs/` dizini - Opsiyonel):**
    *   `dist-sw-conf.ios`, `edge-sw-conf.ios`, `router-conf.ios`: GNS3'teki Cisco anahtar ve router'ları için IOS konfigürasyon şablonlarıdır. `edge-sw-conf.ios` dosyasında Go şablon sentaksı (`{{.HOSTNAME}}`, `{{.IP_ADDRESS}}`) kullanılarak dinamik değerler atanır.
    *   `sw-list.json`: GNS3 projenizdeki switch ve router'ların konsol portları gibi bilgilerini içeren bir JSON dosyasıdır. `edge-configurator.go` aracı bu dosyayı kullanarak cihazlara bağlanır.

## Kullanım

### LAB Ortamını Başlatma

Proje ana dizininde (`l2-observer-lab/`) aşağıdaki komutu çalıştırarak tüm Docker servislerini arka planda başlatın:

```bash
docker-compose up -d --build
```
`--build` parametresi, Go keşif servisi için Docker imajını yeniden derler.

### Doğrulama Adımları

Tüm servisler başladıktan sonra, ortamın düzgün çalıştığını doğrulamak için aşağıdaki adımları uygulayın:

1.  **Keşif Dosyasını Kontrol Etme:**
    *   `discovery/targets/tüm_switchler.json` dosyasının oluştuğunu ve içinde simüle edilmiş switch'lerin IP adreslerinin ve bölgelerinin (`A_Veri_Merkezi`, `B_Veri_Merkezi` vb.) listelendiğini kontrol edin.

2.  **Leaf Prometheus Sunucularını Kontrol Etme:**
    *   `docker-compose ps` komutunu çalıştırarak `leaf-a` servisinin dışarıya açılan portunu bulun (örneğin `0.0.0.0:XXXX->9090/tcp`).
    *   Tarayıcınızda `http://localhost:XXXX` adresini (bulduğunuz portu kullanarak) açın.
    *   "Status" -> "Targets" menüsüne gidin. Burada sadece `A_Veri_Merkezi` bölgesine ait IP adreslerini (örneğin `172.20.0.10`, `172.20.0.11`) görmelisiniz.

3.  **Global Prometheus Sunucusunu Kontrol Etme:**
    *   Tarayıcınızda `http://localhost:9090` adresini açın.
    *   "Status" -> "Targets" menüsüne gidin. Burada tüm Leaf Prometheus sunucularını (`leaf-a:9090`, `leaf-b:9090` vb.) hedeflendiğini görmelisiniz.

4.  **Grafana'yı Kontrol Etme:**
    *   Tarayıcınızda `http://localhost:3000` adresini açın.
    *   Varsayılan kullanıcı adı/şifre `admin/admin` ile giriş yapın. Şifrenizi değiştirmeniz istenebilir.
    *   Yeni bir Prometheus veri kaynağı ekleyin:
        *   **Name:** `Global Prometheus`
        *   **URL:** `http://global-prometheus:9090` (Bu URL Docker ağı içindeki servis adı ve portudur)
        *   Save & Test yapın. Bağlantının başarılı olduğunu doğrulayın.
    *   Ardından, `Explore` (Keşfet) veya yeni bir pano oluşturarak `up` gibi temel bir Prometheus sorgusu ile verilerin geldiğini kontrol edin.

### Ağ Cihazı Konfigürasyon Aracı Kullanımı (Opsiyonel)

Bu araç, bir GNS3 ortamında çalışan gerçek veya sanal Cisco anahtarlarını otomatik olarak yapılandırmak için kullanılır.

1.  **GNS3 Projenizi Başlatın:** GNS3'ü açın ve `prom-snmp-exporter` projesindeki tüm ağ cihazlarını (`DIST-SW1`, `DIST-SW2`, `SW-A0` vb.) başlatın.
2.  **`sw-list.json` Dosyasını Güncelleyin:** GNS3 projenizin dışa aktarılmış `sw-list.json` dosyasını `network-confs/sw-list.json` konumuna kopyalayın. Bu dosya, GNS3'teki cihazların telnet konsol portlarını içerir.
3.  **Routerlara Statik Yönlendirme Ekleyin:** Linux cihazınızda aşağıdaki komutu çalıştırın. Bu, Docker ortamı ile GNS3 arasındaki iletişimi sağlar.
    ```bash
    sudo ip route add 10.0.0.0/24 via 192.168.122.1 dev virbr0
    ```
4.  **Konfigürasyon Aracını Çalıştırın:** `network-confs` dizinine gidin ve aracı çalıştırın:
    ```bash
    cd network-confs
    go run edge-configurator.go
    ```
    Bu araç, `sw-list.json` dosyasındaki telnet konsol portlarını kullanarak SW-X0/X1/X2/X3/X4 formatındaki cihazlara otomatik olarak `edge-sw-conf.ios` şablonunu uygulayacaktır. Komut çıktısını takip ederek konfigürasyon sürecini görebilirsiniz. DIST-SW'ler için `dist-sw-conf.ios` ve Router'lar için `router-conf.ios` dosyalarını manuel olarak uygulamanız gerekebilir.

## Servis Olarak Çalıştırma

Bu proje, bir sunucu ortamında kalıcı olarak çalıştırılmak üzere Docker Compose ile dağıtılmaya uygun bir yapıya sahiptir.

*   **Servisleri Başlatma (Arka Planda):**
    ```bash
    docker-compose up -d
    ```
*   **Tüm Servisleri Durdurma:**
    ```bash
    docker-compose down
    ```
*   **Servisleri Yeniden Başlatma:**
    ```bash
    docker-compose restart
    ```
*   **Servislerin Durumunu Görüntüleme:**
    ```bash
    docker-compose ps
    ```
*   **Servis Loglarını Görüntüleme:**
    Belirli bir servisin çıktılarını takip etmek için:
    ```bash
    docker-compose logs -f <servis_adı>
    # Örneğin: docker-compose logs -f discovery-service
    # Tüm servislerin loglarını görmek için: docker-compose logs -f
    ```

## API Referansı

Bu proje, harici bir RESTful API sunmamaktadır. Temel bileşenler aşağıdaki gibi çalışır:

*   **Prometheus:** Metrikleri standart Prometheus scraping mekanizmaları aracılığıyla `/metrics` endpoint'inde sunar (örneğin, `http://localhost:9090/metrics` veya `http://snmp-exporter:9116/metrics`).
*   **SNMP Exporter:** SNMP metriklerini, belirtilen hedeflerden toplayarak kendi `/metrics` endpoint'inde Prometheus'a sunar.
*   **Keşif Servisi:** Ağ cihazı hedeflerini dahili olarak bir JSON dosyasına (`discovery/targets/tüm_switchler.json`) yazar ve bu dosya Prometheus tarafından okunur. Harici bir HTTP API sağlamaz.

## Lisans

Bu proje **MIT Lisansı** altında lisanslanmıştır. Daha fazla bilgi için `LICENSE` dosyasına bakınız.
---