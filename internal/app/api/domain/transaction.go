//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../test/mock/$GOPACKAGE/$GOFILE
package domain

import (
	"context"
)

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) error) error
}
