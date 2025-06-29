package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"os"
	"strconv"
	"strings"
)

type nodeConf struct {
	HOSTNAME     string `json:"name"`         // SW üzerine verilmiş liste adı. Hostname ile aynı.
	CONSOLE_TYPE string `json:"console_type"` // Telnet ise "telnet" yazar. Kontrol için.
	PORT         int    `json:"console"`      // Konsol portu buradan telnet portunu alacağız.
}

type nodeSettings struct {
	HOSTNAME   string
	IP_ADDRESS string
}

func getNodes() []nodeConf {
	// Dosyadan SW listesini çek
	nodesFiles, err := os.ReadFile("sw-list.json")
	if err != nil {
		fmt.Println("Dosya okunamadı:", err)
		return nil
	}
	// Tüm düğümleri okuma
	var allNodes []any
	err = json.Unmarshal(nodesFiles, &allNodes)
	if err != nil {
		fmt.Println("JSON parse edilemedi:", err)
		return nil
	}
	// Okunan düğümleri filtreleme
	var nodes []nodeConf
	for i, anyNode := range allNodes {
		// Geri String çevir
		nodeBytes, err := json.Marshal(anyNode)
		if err != nil {
			fmt.Printf("  -> Eleman #%d marshal edilemedi, atlanıyor.\n", i+1)
			continue
		}
		var node nodeConf
		if err := json.Unmarshal(nodeBytes, &node); err != nil {
			fmt.Printf("  -> Eleman #%d parse edilemedi, atlanıyor.\n", i+1)
			continue
		} else if node.CONSOLE_TYPE != "telnet" {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func telnetConfigWriter(node nodeConf, config bytes.Buffer) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(node.PORT))
	if err != nil {
		fmt.Println("Hata! Sunucuya bağlanılamadı:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Printf("-> %s:%d adresine bağlandı. Sunucu bekleniyor...\n", "127.0.0.1", node.PORT)

	fmt.Println("-> Sunucuya ilk 'Enter' sinyali gönderiliyor...")
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		fmt.Println("İlk 'Enter' gönderilirken hata:", err)
		return
	}
	fmt.Println("-> Sunucudan yanıt bekleniyor...")

	for {
		// Sunucudan gelen satırı oku
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nSunucuyla bağlantı koptu: %v\n", err)
			return // Hata durumunda döngüden çık
		}

		// Gelen satırı ekrana yazarak ne olduğunu gör
		fmt.Print(line)

		// Gelen satırın içeriğine göre karar ver
		// NOT: "else if" yerine ayrı "if" blokları kullanmak bazen daha esnek olabilir
		// ama bu yapı için "else if" daha uygun.
		if strings.Contains(line, "(config)#") {
			fmt.Println("-> Config modda, ayarlar gönderiliyor...")
			_, err = conn.Write(config.Bytes())
			if err != nil {
				fmt.Println("Konfigürasyon gönderilirken hata:", err)
				return
			}
			// Bütün konfigürasyon gönderildikten sonra çıkış komutunu gönderelim
			// ve döngüyü sonlandıralım.
			fmt.Println("-> Ayarlar gönderildi. Oturum kapatılıyor.")
			conn.Write([]byte("end\n"))  // Config moddan çık
			conn.Write([]byte("wr\n"))   // Yapılan ayarları yaz
			conn.Write([]byte("exit\n")) // Cihazdan çık
			break                        // İşimiz bitti, döngüyü sonlandır.

		} else if strings.Contains(line, "#") {
			fmt.Println("-> Privileged modda, config moda geçiliyor...")
			_, err = conn.Write([]byte("configure terminal\n")) // "conf t" yerine tam komut daha güvenilir
			if err != nil {
				fmt.Println("Komut gönderilirken hata:", err)
				return
			}

		} else if strings.Contains(line, ">") {
			fmt.Println("-> User modda, privileged moda geçiliyor...")
			_, err = conn.Write([]byte("enable\n"))
			if err != nil {
				fmt.Println("Komut gönderilirken hata:", err)
				return
			}

			// "Press RETURN" veya benzeri bir başlangıç mesajı varsa diye...
		} else if strings.Contains(line, "Press RETURN") {
			_, err = conn.Write([]byte("\n"))
			if err != nil {
				fmt.Println("Komut gönderilirken hata:", err)
				return
			}
		}
	}
	fmt.Println("-> Oturum başarıyla sonlandırıldı.")
}

func main() {
	// getNodes fonksiyonunu da daha verimli haliyle değiştirmeyi unutma
	nodes := getNodes()
	confTemplate, err := template.ParseFiles("edge-sw-conf.ios")
	if err != nil {
		panic(err)
	}
	var config bytes.Buffer
	for _, node := range nodes {
		// Döngü her başladığında buffer'ı sıfırla!
		config.Reset()

		if node.HOSTNAME[0:2] == "SW" {
			areaNumber := []rune(node.HOSTNAME[3:4])[0] - []rune("A")[0] + 1
			fmt.Printf("Hostname To IP\nHostname: %s\nIP: 192.168.100.%d%s\n", node.HOSTNAME, areaNumber, node.HOSTNAME[4:5])
			fmt.Printf("Ayar basılıyor... / PORT: %d\n", node.PORT)

			err := confTemplate.Execute(&config, nodeSettings{
				HOSTNAME:   node.HOSTNAME,
				IP_ADDRESS: fmt.Sprintf("192.168.100.%d%s", areaNumber, node.HOSTNAME[4:5]),
			})

			if err != nil {
				fmt.Println("Şablon hatası:", err)
				continue
			}

			telnetConfigWriter(node, config)
		}
	}
}
