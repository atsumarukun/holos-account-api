package transformer

import (
	"github.com/atsumarukun/holos-account-api/internal/app/api/domain/entity"
	"github.com/atsumarukun/holos-account-api/internal/app/api/infrastructure/model"
	"github.com/google/uuid"
)

func ToAccountModel(account *entity.Account) *model.AccountModel {
	return &model.AccountModel{
		ID:       account.ID.String(),
		Name:     account.Name,
		Password: account.Password,
	}
}

func ToAccountEntity(account *model.AccountModel) (*entity.Account, error) {
	id, err := uuid.Parse(account.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Account{
		ID:       id,
		Name:     account.Name,
		Password: account.Password,
	}, nil
}
