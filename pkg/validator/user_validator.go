package validator

import (
	"context"
	"errors"
	"strings"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/user"
)

func ValidateEmailUniqueness(ctx context.Context, db *ent.Client, email string) error {
	exists, err := db.User.Query().
		Where(user.EmailEQ(strings.ToLower(email))).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}
	return nil
}

func ValidatePhoneUniqueness(ctx context.Context, db *ent.Client, phone string) error {
	if phone == "" {
		return nil
	}
	exists, err := db.User.Query().
		Where(
			user.PhoneNumberEQ(phone),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("phone number already exists")
	}
	return nil
}
