package application

type PrinterEngine interface {
    AlignLeft()
    AlignCenter()
    AlignRight()
    BoldOn()
    BoldOff()
    UnderlineOn()
    UnderlineOff()
    Size(width, height int)
    SizeReset()
    NewLine()
    NewLines(n int)
    AppendRaw(s string)
    Cut()
    PrintAndFlush() error
    PrintImageURL(url string) error
    PrintImageURLScaled(url string, maxWidth int) error
}