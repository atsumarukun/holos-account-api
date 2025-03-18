//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../../../../../test/mock/domain/repository/$GOPACKAGE/$GOFILE
package helper

import "context"

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) error) error
}
