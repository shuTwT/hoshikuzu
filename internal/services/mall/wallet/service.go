package wallet

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/wallet"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type WalletService interface {
	QueryWallet(c *fiber.Ctx, userId int) (*ent.Wallet, error)
	QueryWalletPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Wallet, error)
	CreateWallet(c context.Context, userId int) (*ent.Wallet, error)
	UpdateWallet(c context.Context, walletId int, walletData model.WalletUpdateReq) (*ent.Wallet, error)
	UpdateWalletBalance(c context.Context, walletId int, balance int) (*ent.Wallet, error)
	FreezeWallet(c context.Context, walletId int, amount int) (*ent.Wallet, error)
	UnfreezeWallet(c context.Context, walletId int, amount int) (*ent.Wallet, error)
}

type WalletServiceImpl struct {
	client *ent.Client
}

func NewWalletServiceImpl(client *ent.Client) *WalletServiceImpl {
	return &WalletServiceImpl{client: client}
}

func (s *WalletServiceImpl) QueryWallet(c *fiber.Ctx, userId int) (*ent.Wallet, error) {
	w, err := s.client.Wallet.Query().
		Where(wallet.UserIDEQ(userId)).
		Only(c.Context())
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletServiceImpl) QueryWalletPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Wallet, error) {
	count, err := s.client.Wallet.Query().Count(c.UserContext())
	if err != nil {
		return 0, nil, err
	}

	wallets, err := s.client.Wallet.Query().
		Order(ent.Desc(wallet.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		return 0, nil, err
	}

	return count, wallets, nil
}

func (s *WalletServiceImpl) CreateWallet(c context.Context, userId int) (*ent.Wallet, error) {
	w, err := s.client.Wallet.Create().
		SetUserID(userId).
		Save(c)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletServiceImpl) UpdateWallet(c context.Context, id int, updateReq model.WalletUpdateReq) (*ent.Wallet, error) {
	update := s.client.Wallet.UpdateOneID(id)

	if updateReq.Password != "" {
		update.SetPassword(updateReq.Password)
	}

	if updateReq.Remark != "" {
		update.SetRemark(updateReq.Remark)
	}

	if updateReq.Active != nil {
		update.SetActive(*updateReq.Active)
	}

	updatedWallet, err := update.Save(c)
	if err != nil {
		return nil, err
	}
	return updatedWallet, nil
}

func (s *WalletServiceImpl) UpdateWalletBalance(c context.Context, walletId int, balance int) (*ent.Wallet, error) {

	w, err := s.client.Wallet.UpdateOneID(walletId).
		SetBalance(balance).
		Save(c)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletServiceImpl) FreezeWallet(c context.Context, walletId int, amount int) (*ent.Wallet, error) {

	w, err := s.client.Wallet.UpdateOneID(walletId).
		AddFrozenAmount(amount).
		AddBalance(-amount).
		Save(c)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletServiceImpl) UnfreezeWallet(c context.Context, walletId int, amount int) (*ent.Wallet, error) {
	w, err := s.client.Wallet.UpdateOneID(walletId).
		AddFrozenAmount(-amount).
		AddBalance(amount).
		Save(c)
	if err != nil {
		return nil, err
	}
	return w, nil
}
