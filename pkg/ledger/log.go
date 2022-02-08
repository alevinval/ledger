package ledger

import (
	"github.com/alevinval/ledger/internal/log"
	"go.uber.org/zap"
)

func SetZapLogger(l *zap.Logger) {
	log.SetZapLogger(l)
}
