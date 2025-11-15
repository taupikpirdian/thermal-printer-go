package domain

import (
    "fmt"
    "regexp"
    "strings"
    "unicode/utf8"
)

func RepeatRune(r rune, n int) string {
    if n <= 0 {
        return ""
    }
    return strings.Repeat(string(r), n)
}

func FormatNumberID(n int64) string {
    s := fmt.Sprintf("%d", n)
    neg := false
    if strings.HasPrefix(s, "-") {
        neg = true
        s = s[1:]
    }
    var out []string
    for len(s) > 3 {
        out = append([]string{s[len(s)-3:]}, out...)
        s = s[:len(s)-3]
    }
    if s != "" {
        out = append([]string{s}, out...)
    }
    res := strings.Join(out, ".")
    if neg {
        res = "-" + res
    }
    return res
}

func wrapTextToWidth(s string, width int) []string {
    s = strings.ReplaceAll(s, "\r\n", "\n")
    lines := []string{}
    for _, raw := range strings.Split(s, "\n") {
        line := raw
        for line != "" {
            if runeLen(line) <= width {
                lines = append(lines, line)
                break
            }
            prefix := takeRunes(line, width)
            lines = append(lines, prefix)
            line = line[len(prefix):]
        }
        if len(raw) == 0 {
            lines = append(lines, "")
        }
    }
    return lines
}

func runeLen(s string) int {
    return utf8.RuneCountInString(s)
}

func takeRunes(s string, n int) string {
    if n <= 0 {
        return ""
    }
    var b strings.Builder
    i := 0
    for _, r := range s {
        if i >= n {
            break
        }
        b.WriteRune(r)
        i++
    }
    return b.String()
}

func PrintInvoiceItem(w LineWriter, paperWidth int, left, right string) {
    re := regexp.MustCompile(`\r`)
    left = re.ReplaceAllString(left, "")
    right = re.ReplaceAllString(right, "")
    left = strings.TrimSpace(left)
    right = strings.TrimSpace(right)
    space := paperWidth - runeLen(left) - runeLen(right)
    if space >= 0 {
        line := left + strings.Repeat(" ", space) + right
        w.AppendRaw(line)
        w.NewLine()
        return
    }
    w.AppendRaw(left)
    w.NewLine()
    for _, l := range wrapTextToWidth(right, paperWidth) {
        w.AppendRaw(l)
        w.NewLine()
    }
}

func PrintLeftAlignedText(w LineWriter, paperWidth int, s string) {
    for _, l := range wrapTextToWidth(s, paperWidth) {
        w.AppendRaw(strings.TrimSpace(l))
        w.NewLine()
    }
}

func ProcessName(s string) string {
    s = strings.TrimSpace(s)
    if runeLen(s) > 20 {
        return takeRunes(s, 20)
    }
    return s
}