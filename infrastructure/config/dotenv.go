package config

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "os"
    "strings"
    "path/filepath"
)

func LoadDotEnv(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    s := bufio.NewScanner(f)
    for s.Scan() {
        line := strings.TrimSpace(s.Text())
        if line == "" || strings.HasPrefix(line, "#") { continue }
        i := strings.IndexByte(line, '=')
        if i <= 0 { continue }
        key := strings.TrimSpace(line[:i])
        val := strings.TrimSpace(line[i+1:])
        if len(val) >= 2 {
            if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
                val = val[1:len(val)-1]
            }
        }
        _ = os.Setenv(key, val)
    }
    return nil
}

func LoadDotEnvIfExists(path string) error {
    if strings.TrimSpace(path) == "" { return nil }
    if _, err := os.Stat(path); err != nil { return nil }
    return LoadDotEnv(path)
}

var EmbeddedEnvBase64 string

func LoadEmbeddedEnv() {
    s := strings.TrimSpace(EmbeddedEnvBase64)
    if s == "" { return }
    b, err := base64.StdEncoding.DecodeString(s)
    if err != nil { return }
    var m map[string]string
    if json.Unmarshal(b, &m) != nil { return }
    for k, v := range m {
        _ = os.Setenv(k, v)
    }
}

func LoadDefaultEnv() {
    LoadEmbeddedEnv()
    _ = LoadDotEnvIfExists(".env")
    exe, err := os.Executable()
    if err == nil {
        _ = LoadDotEnvIfExists(filepath.Join(filepath.Dir(exe), ".env"))
    }
    appdata := os.Getenv("APPDATA")
    if appdata != "" {
        _ = LoadDotEnvIfExists(filepath.Join(appdata, "MidasPrinter", ".env"))
    }
    programData := os.Getenv("ProgramData")
    if programData != "" {
        _ = LoadDotEnvIfExists(filepath.Join(programData, "MidasPrinter", ".env"))
    }
}