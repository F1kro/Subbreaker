package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Banner
func tampilkanBanner() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║     SUBDOMAIN ENUMERATOR & PORT SCAN     ║")
	fmt.Println("║             by : Fiqro Najiah            ║")
	fmt.Println("║         Github : github.com/F1kro        ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}

func tampilkanHelp() {
	fmt.Println("Penggunaan:")
	fmt.Println("  -d <domain>        : domain target (wajib)")
	fmt.Println("  -o <namafile.txt>  : simpan hasil ke file (opsional)")
	fmt.Println("  -h                 : tampilkan bantuan")
}

func validasiOutputFile(nama string) error {
	if !strings.HasSuffix(nama, ".txt") {
		return fmt.Errorf("format file output harus .txt")
	}
	dir := filepath.Dir(nama)
	if _, err := os.Stat(dir); os.IsNotExist(err) && dir != "." {
		return fmt.Errorf("direktori tidak ditemukan: %s", dir)
	}
	return nil
}

//  crt.sh
func ambilDariCRTSh(domain string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var hasil []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&hasil); err != nil {
		return nil
	}

	set := make(map[string]bool)
	for _, entry := range hasil {
		if name, ok := entry["name_value"].(string); ok {
			for _, sub := range strings.Split(name, "\n") {
				if strings.HasSuffix(sub, domain) {
					set[sub] = true
				}
			}
		}
	}

	var unik []string
	for sub := range set {
		unik = append(unik, sub)
	}
	return unik
}

// findsubdomains
func ambilDariFindSubdomains(domain string) []string {
	url := fmt.Sprintf("https://findsubdomains.com/api/subdomains/%s", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Subdomains []string `json:"subdomains"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}

	var hasil []string
	for _, sub := range data.Subdomains {
		hasil = append(hasil, fmt.Sprintf("%s.%s", sub, domain))
	}
	return hasil
}

func resolusiIP(host string) string {
	ip, err := net.LookupIP(host)
	if err != nil || len(ip) == 0 {
		return "IP tidak ditemukan"
	}
	return ip[0].String()
}

func scanPort(host string, portList []int) []int {
	var terbuka []int
	for _, port := range portList {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
		if err == nil {
			terbuka = append(terbuka, port)
			conn.Close()
		}
	}
	return terbuka
}

func simpanKeFile(nama string, data []string) error {
	file, err := os.Create(nama)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, baris := range data {
		fmt.Fprintln(writer, baris)
	}
	return writer.Flush()
}

func main() {
	domain := flag.String("d", "", "domain target")
	output := flag.String("o", "", "nama file output")
	help := flag.Bool("h", false, "tampilkan bantuan")
	flag.Parse()

	if *help || *domain == "" {
		tampilkanBanner()
		tampilkanHelp()
		return
	}

	tampilkanBanner()

	if *output != "" {
		if err := validasiOutputFile(*output); err != nil {
			fmt.Println("❌ Error:", err)
			os.Exit(1)
		}
		if _, err := os.Stat(*output); err == nil {
			fmt.Printf("❌ File '%s' sudah ada. Ganti nama file atau hapus file lama dulu.\n", *output)
			os.Exit(1)
		}
	}

	fmt.Printf("Mengambil subdomain dari %s ...\n\n", *domain)

	subdomainSet := make(map[string]bool)

	sumber1 := ambilDariCRTSh(*domain)
	sumber2 := ambilDariFindSubdomains(*domain)

	for _, sub := range append(sumber1, sumber2...) {
		subdomainSet[sub] = true
	}

	if len(subdomainSet) == 0 {
		fmt.Println("❌ Tidak ditemukan subdomain. Periksa kembali domain yang dimasukkan / Periksa koneksi internet anda!.")
		os.Exit(1)
	}

	var subdomainList []string
	for sub := range subdomainSet {
		subdomainList = append(subdomainList, sub)
	}
	sort.Strings(subdomainList)

	portUmum := []int{80, 443, 8080, 22, 21, 25, 3306}
	var hasil []string

	for i, sub := range subdomainList {
		ip := resolusiIP(sub)
		port := scanPort(sub, portUmum)
		portStr := "tidak ada"
		if len(port) > 0 {
			portStr = fmt.Sprintf("%v", port)
		}
		baris := fmt.Sprintf("%d. %s [%s] → http://%s | Port terbuka: %s", i+1, sub, ip, sub, portStr)
		fmt.Println(baris)
		hasil = append(hasil, baris)
	}

	if *output != "" {
		if len(hasil) == 0 {
			fmt.Println("❌ Tidak ada data yang disimpan karena hasil kosong.")
			os.Exit(1)
		}
		if err := simpanKeFile(*output, hasil); err != nil {
			fmt.Println("❌ Gagal menyimpan ke file:", err)
		} else {
			fmt.Println("\n✅ Hasil disimpan ke:", *output)
		}
	}
}
