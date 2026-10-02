package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var (
	blacklist = make(map[string]bool)
	mutex     sync.Mutex
)

func main() {
	// Modbus TCP standart portu 502'dir. 
	port := "502"
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("[-] ModbusTrap başlatılamadı (502 portu için yönetici izni gerekebilir):", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("[+] ModbusTrap AKTİF | ICS/SCADA Savunma Modülü")
	fmt.Printf("[+] Kritik Altyapı Simülasyonu: Soğutma Pompası-A (Port %s)\n", port)
	fmt.Println("---------------------------------------------------")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleIndustrialThreat(conn)
	}
}

func handleIndustrialThreat(conn net.Conn) {
	ipStr := conn.RemoteAddr().String()
	victimIP, _, _ := net.SplitHostPort(ipStr)

	// Zero-Trust IPS Kontrolü
	mutex.Lock()
	isBanned := blacklist[victimIP]
	mutex.Unlock()

	if isBanned {
		fmt.Printf("[🛡️ BLOKE] Yasaklı IP endüstriyel ağdan atıldı: %s\n", victimIP)
		conn.Close()
		return
	}

	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	
	if err != nil || n == 0 {
		return // Sessiz taramaları geçiyoruz, hedefimiz veri gönderen sabotajcılar
	}

	fmt.Printf("\n[🚨 KRİTİK ALARM] Endüstriyel Ağa Sızma Girişimi: %s\n", victimIP)
	
	payload := buffer[:n]
	attackType := "Bilinmeyen TCP Verisi"

	// Modbus TCP paket analizi (Stuxnet benzeri keşifleri yakalamak için)
	// Modbus TCP paketinin 2. ve 3. byte'ı Protocol ID'dir ve her zaman 0x00 0x00 olmalıdır.
	if n >= 8 && payload[2] == 0x00 && payload[3] == 0x00 {
		functionCode := payload[7] // 7. byte fonksiyon kodunu tutar (Okuma/Yazma)
		
		if functionCode == 0x01 || functionCode == 0x02 || functionCode == 0x03 || functionCode == 0x04 {
			attackType = fmt.Sprintf("SCADA Keşif Taraması (Modbus Function 0x%02x - Veri Okuma İsteği)", functionCode)
		} else if functionCode == 0x05 || functionCode == 0x06 || functionCode == 0x0F || functionCode == 0x10 {
			attackType = fmt.Sprintf("ENDÜSTRİYEL SABOTAJ! (Modbus Function 0x%02x - Vana/Pompa Durumu Değiştirme İsteği)", functionCode)
		} else {
			attackType = "Şüpheli Modbus Trafiği"
		}
	}

	fmt.Printf("[!] Sınıflandırma: %s\n", attackType)
	fmt.Printf("[!] Gelen Ham Veri (Hex): %s\n", hex.EncodeToString(payload))

	// Anında Karantina
	mutex.Lock()
	blacklist[victimIP] = true
	mutex.Unlock()
	fmt.Printf("[🛑 İZOLASYON] %s IP adresi SCADA ağından tamamen kilitlendi!\n", victimIP)
	fmt.Println("---------------------------------------------------")
}