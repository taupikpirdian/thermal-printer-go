package escpos

import (
    "bytes"
    "strings"
    "github.com/alexbrainman/printer"
)

type EscposPrinter struct {
    p   *printer.Printer
    buf bytes.Buffer
}

func NewEscposPrinter(p *printer.Printer) *EscposPrinter {
    return &EscposPrinter{p: p}
}

func (e *EscposPrinter) write(b []byte) { e.buf.Write(b) }
func (e *EscposPrinter) writeString(s string) {
    s = strings.ReplaceAll(s, "\r\n", "\n")
    e.buf.WriteString(s)
}

func (e *EscposPrinter) flush() error {
    _, err := e.p.Write(e.buf.Bytes())
    e.buf.Reset()
    return err
}

func (e *EscposPrinter) NewLine()                 { e.write([]byte("\n")) }
func (e *EscposPrinter) NewLines(n int)            { for i := 0; i < n; i++ { e.NewLine() } }
func (e *EscposPrinter) AlignLeft()                { e.write([]byte{0x1B, 0x61, 0x00}) }
func (e *EscposPrinter) AlignCenter()              { e.write([]byte{0x1B, 0x61, 0x01}) }
func (e *EscposPrinter) AlignRight()               { e.write([]byte{0x1B, 0x61, 0x02}) }
func (e *EscposPrinter) BoldOn()                   { e.write([]byte{0x1B, 0x45, 0x01}) }
func (e *EscposPrinter) BoldOff()                  { e.write([]byte{0x1B, 0x45, 0x00}) }
func (e *EscposPrinter) UnderlineOn()              { e.write([]byte{0x1B, 0x2D, 0x01}) }
func (e *EscposPrinter) UnderlineOff()             { e.write([]byte{0x1B, 0x2D, 0x00}) }
func (e *EscposPrinter) Size(width, height int)    { if width < 1 { width = 1 }; if height < 1 { height = 1 }; if width > 8 { width = 8 }; if height > 8 { height = 8 }; n := byte(((width - 1) << 4) | (height - 1)); e.write([]byte{0x1D, 0x21, n}) }
func (e *EscposPrinter) SizeReset()                { e.write([]byte{0x1D, 0x21, 0x00}) }
func (e *EscposPrinter) Cut()                      { e.write([]byte("\n\n\n")); e.write([]byte{0x1D, 0x56, 0x01}) }
func (e *EscposPrinter) AppendRaw(s string)        { e.writeString(s) }
func (e *EscposPrinter) PrintAndFlush() error      { return e.flush() }

func (e *EscposPrinter) PrintImageURL(url string) error {
    if strings.TrimSpace(url) == "" { return nil }
    img, err := loadImage(url)
    if err != nil { return err }
    e.write([]byte{0x1B, 0x40})
    e.write([]byte{0x1B, 0x61, 0x02})
    e.write(convertToBitImage(img))
    e.write([]byte{0x1B, 0x61, 0x00})
    return nil
}

func (e *EscposPrinter) PrintImageURLScaled(url string, maxWidth int) error {
    if strings.TrimSpace(url) == "" { return nil }
    img, err := loadImage(url)
    if err != nil { return err }
    scaled := resizeToWidth(img, maxWidth)
    e.write([]byte{0x1B, 0x40})
    e.write([]byte{0x1B, 0x61, 0x02})
    e.write(convertToBitImage(scaled))
    e.write([]byte{0x1B, 0x61, 0x00})
    return nil
}

func (e *EscposPrinter) PrintImageFile(path string) error {
    if strings.TrimSpace(path) == "" { return nil }
    img, err := loadLocalImage(path)
    if err != nil { return err }
    e.write([]byte{0x1B, 0x40})
    e.write([]byte{0x1B, 0x61, 0x02})
    e.write(convertToBitImage(img))
    e.write([]byte{0x1B, 0x61, 0x00})
    return nil
}