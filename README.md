# Go Thermal Printer (ESC/POS) — EPSON TM-U220

Aplikasi Go untuk mencetak struk dan gambar ke printer thermal menggunakan ESC/POS. Arsitektur mengikuti Domain-Driven Design (DDD) dan Clean Architecture. Input data diterima via MQTT, konfigurasi dari environment dan fallback ke `setting.xml`.

Telah diuji pada printer EPSON TM-U220 di Windows.

## Fitur
- Cetak teks dan gambar ESC/POS (mode 8-dot) yang kompatibel dengan TM-U220
- Scaling dan padding gambar agar proporsional di kertas kecil
- Subscribe MQTT untuk menerima payload JSON dan langsung cetak
- Loader `.env` sederhana tanpa dependency eksternal
- Fallback konfigurasi dari `setting.xml` (kompatibel dengan aplikasi lama)
- Logging terstruktur saat startup dan saat pesan MQTT diterima

## Struktur Proyek
```
cmd/printer/            Entrypoint aplikasi
application/            Use case (PrintInvoiceUseCase)
domain/                 Entitas, utilitas bisnis (format, wrapping)
infrastructure/escpos/  Encoder dan driver ESC/POS (Windows printer)
infrastructure/mqttclient/ MQTT client (Paho)
infrastructure/config/   Loader .env dan settings.xml
```

## Prasyarat
- Go `1.24` (mengacu `go.mod`)
- Windows dengan printer terpasang di Control Panel (uji: EPSON TM-U220)
- Broker MQTT yang dapat diakses

## Konfigurasi
Gunakan `.env` di root repository. Lihat contoh lengkap di `.env.example`.

Variabel penting:
- `MQTT_BROKER_URL` atau `MQTT_HOST` + `MQTT_PORT`
- `MQTT_USERNAME` + `MQTT_PASSWORD` atau `MQTT_AUTH_USERNAME` + `MQTT_AUTH_PASSWORD`
- `MQTT_TOPIC_PREFIX` dan `ID_MERCHANT` atau langsung `MQTT_TOPIC`
- `PRINTER_NAME` contoh: `EPSON TM-U220 Receipt`
- `PAPER_WIDTH` default `42`
- `OFFSITE` padding horizontal default `8`

Optional `setting.xml` di working directory:
```
<root>
  <id_merchant>...</id_merchant>
  <name_merchant>...</name_merchant>
  <printer_name>EPSON TM-U220 Receipt</printer_name>
  <size>42</size>
  <offsite>8</offsite>
</root>
```
Nilai env akan mengoverride file ini.

## Menjalankan dari Source
Pastikan `.env` sudah tersedia di root project, lalu jalankan:
```
go run ./cmd/printer
```

Jika berhasil, terminal akan menampilkan:
```
Client ready | Broker: tcp://<host>:<port> | Topic: <topic> | Printer: <printer>
```

Catatan:
- Printing langsung ke Windows spooler hanya berjalan di Windows.
- Di macOS/Linux, aplikasi bisa di-build untuk validasi compile, tetapi saat menerima pesan MQTT akan gagal membuka printer dengan pesan bahwa raw Windows printer tidak didukung.

## Build Executable

### Build di Windows
Jalankan dari root project:
```
go build -o valas-print-midas.exe ./cmd/printer
```

Lalu jalankan:
```
.\valas-print-midas.exe
```

### Build `.exe` dari macOS/Linux
Gunakan cross compile:
```
GOOS=windows GOARCH=amd64 go build -o valas-print-midas.exe ./cmd/printer
```

Jika environment membatasi lokasi Go build cache, arahkan cache ke folder sementara:
```
GOCACHE=/tmp/thermal-printer-go-cache GOOS=windows GOARCH=amd64 go build -o valas-print-midas.exe ./cmd/printer
```

### Build untuk Validasi Lokal
Untuk memastikan semua package compile di OS saat ini:
```
go test ./...
go build -o /tmp/thermal-printer-go-printer ./cmd/printer
```

## Menjalankan Executable
Letakkan `.env` di salah satu lokasi berikut:
- Folder tempat executable dijalankan
- Folder yang sama dengan file `.exe`
- `%APPDATA%\MidasPrinter\.env`
- `%ProgramData%\MidasPrinter\.env`

Jalankan executable. Saat berhasil inisialisasi, terminal menampilkan:
```
Client ready | Broker: tcp://<host>:<port> | Topic: <topic> | Printer: <printer>
```

Saat pesan MQTT diterima, terlihat log:
```
MQTT received | topic=<topic> | bytes=<len>
MQTT payload: { ...json... }
Parsed payload | order=<no_order> | merchant=<merchant_name>
Printer opened | name=<printer_name>
StartRawDocument OK
Image printed, performing cut
Flush OK
Pencetakan Selesai
```

## Format Payload
Aplikasi mengharapkan JSON seperti berikut (dipersingkat):
```
{
  "merchant_name": "Midas xchange",
  "no_order": "FB202511147464",
  "tanggal": "2025-11-14",
  "detail_order": [ ... ],
  "logo": "https://.../logo.png",
  "header": { "title": "INVOICE PEMBELIAN", ... },
  "footer": { "note_id": "...", "note_en": "..." }
}
```
Entrypoint saat ini memanggil `Execute` di `cmd/printer/main.go` untuk mencetak struk lengkap.

## Tuning Gambar
- Ubah lebar gambar di `application/print_invoice.go` pada pemanggilan `PrintImageURLScaled(payload.Logo, 200)`
  - Nilai aman untuk TM-U220: `160–220` dot
- Padding kiri/kanan dapat disesuaikan di `infrastructure/escpos/printer.go` melalui `addHorizontalPadding`

## Troubleshooting
- Karakter acak saat cetak gambar: gunakan mode 8-dot `ESC * 1` (sudah default); hindari 24-dot pada TM-U220.
- Gambar terpotong: turunkan `maxWidth` atau tambah padding horizontal (`OFFSITE`).
- `No connection could be made ... [::1]:1883`: pastikan env `MQTT_HOST`/`MQTT_BROKER_URL` terbaca. Gunakan `.env` dan loader sudah aktif.
- `open printer` gagal: pastikan `PRINTER_NAME` sama persis dengan nama di Control Panel.

## Keamanan
- Jangan commit `.env` berisi kredensial. Gunakan `.env.example` sebagai referensi.

## Lisensi
Gunakan sesuai kebutuhan Anda. Jika akan dipublikasikan, tambahkan lisensi sesuai preferensi.
