package config

import (
    "encoding/xml"
    "os"
    "strconv"
)

type Settings struct {
    IDMerchant  string
    NameMerchant string
    PrinterName string
    PaperWidth  int
    Offsite     int
    TopicPrefix string
}

type xmlSettings struct {
    IDMerchant   string `xml:"id_merchant"`
    NameMerchant string `xml:"name_merchant"`
    PrinterName  string `xml:"printer_name"`
    Size         string `xml:"size"`
    Offsite      string `xml:"offsite"`
}

func LoadSettings() (Settings, error) {
    s := Settings{}
    s.IDMerchant = os.Getenv("ID_MERCHANT")
    s.NameMerchant = os.Getenv("NAME_MERCHANT")
    s.PrinterName = os.Getenv("PRINTER_NAME")
    if v := os.Getenv("PAPER_WIDTH"); v != "" {
        if n, err := strconv.Atoi(v); err == nil { s.PaperWidth = n }
    }
    if v := os.Getenv("OFFSITE"); v != "" {
        if n, err := strconv.Atoi(v); err == nil { s.Offsite = n }
    }
    s.TopicPrefix = os.Getenv("MQTT_TOPIC_PREFIX")
    if s.TopicPrefix == "" { s.TopicPrefix = "printer/" }
    if s.PrinterName == "" { s.PrinterName = "EPSON TM-U220 Receipt" }
    if s.PaperWidth == 0 { s.PaperWidth = 42 }

    if s.IDMerchant == "" || s.NameMerchant == "" {
        b, err := os.ReadFile("setting.xml")
        if err == nil {
            var x xmlSettings
            if xml.Unmarshal(b, &x) == nil {
                if s.IDMerchant == "" { s.IDMerchant = x.IDMerchant }
                if s.NameMerchant == "" { s.NameMerchant = x.NameMerchant }
                if s.PrinterName == "" && x.PrinterName != "" { s.PrinterName = x.PrinterName }
                if s.PaperWidth == 0 {
                    if n, err := strconv.Atoi(x.Size); err == nil { s.PaperWidth = n }
                }
                if s.Offsite == 0 {
                    if n, err := strconv.Atoi(x.Offsite); err == nil { s.Offsite = n }
                }
            }
        }
    }
    return s, nil
}