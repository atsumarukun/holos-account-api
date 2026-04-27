package transformer

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/database/model"
)

func ToAccountModel(account *entity.Account) *model.AccountModel {
	if account == nil {
		return nil
	}

	return &model.AccountModel{
		ID:       account.ID,
		Name:     account.Name,
		Password: account.Password,
	}
}

func ToAccountEntity(account *model.AccountModel) *entity.Account {
	if account == nil {
		return nil
	}

	return entity.RestoreAccount(account.ID, account.Name, account.Password)
}
