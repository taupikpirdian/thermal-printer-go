package config

import (
    "bufio"
    "os"
    "strings"
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