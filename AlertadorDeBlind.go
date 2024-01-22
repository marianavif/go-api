package poker

import (
	"fmt"
	"os"
	"time"
)

type AlertadorDeBlind interface {
	AgendarAlertaAs(duracao time.Duration, quantidade int)
}

type AlertadorDeBlindFunc func(duracao time.Duration, quantidade int)

func (a AlertadorDeBlindFunc) AgendarAlertaAs(duracao time.Duration, quantidade int) {
	a(duracao, quantidade)
}

func StdOutAlertador(duracao time.Duration, quantidade int) {
	time.AfterFunc(duracao, func() {
		fmt.Fprintf(os.Stdout, "Blind é agora %d\n", quantidade)
	})
}
