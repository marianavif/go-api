package poker

import (
	"fmt"
	"io"
	"time"
)

type AlertadorDeBlind interface {
	AgendarAlertaPara(duracao time.Duration, quantidade int, para io.Writer)
}

type AlertadorDeBlindFunc func(duracao time.Duration, quantidade int, para io.Writer)

func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantidade int, para io.Writer) {
	a(duracao, quantidade, para)
}

func Alertador(duracao time.Duration, quantidade int, para io.Writer) {
	time.AfterFunc(duracao, func() {
		fmt.Fprintf(para, "Blind é agora %d\n", quantidade)
	})
}
