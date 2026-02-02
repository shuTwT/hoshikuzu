package license

import (
	"context"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	license_ent "github.com/shuTwT/hoshikuzu/ent/license"
)

type LicenseService interface {
	ListLicensePage(ctx context.Context, page, size int) (int, []*ent.License, error)
	QueryLicense(ctx context.Context, id int) (*ent.License, error)
	CreateLicense(ctx context.Context, domain, licenseKey, customerName string, expireDate time.Time) (*ent.License, error)
	UpdateLicense(ctx context.Context, id int, domain, licenseKey, customerName string, expireDate time.Time, status int) (*ent.License, error)
	DeleteLicense(ctx context.Context, id int) error
	VerifyLicense(ctx context.Context, domain string) (*ent.License, error)
}

type LicenseServiceImpl struct {
	client *ent.Client
}

func NewLicenseServiceImpl(client *ent.Client) *LicenseServiceImpl {
	return &LicenseServiceImpl{client: client}
}

func (s *LicenseServiceImpl) ListLicensePage(ctx context.Context, page, size int) (int, []*ent.License, error) {
	count, err := s.client.License.Query().Count(ctx)
	if err != nil {
		return 0, nil, err
	}

	licenses, err := s.client.License.Query().
		Order(ent.Desc(license_ent.FieldID)).
		Offset((page - 1) * size).
		Limit(size).
		All(ctx)
	if err != nil {
		return 0, nil, err
	}

	return count, licenses, nil
}

func (s *LicenseServiceImpl) QueryLicense(ctx context.Context, id int) (*ent.License, error) {
	licenseEntity, err := s.client.License.Query().
		Where(license_ent.ID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return licenseEntity, nil
}

func (s *LicenseServiceImpl) CreateLicense(ctx context.Context, domain, licenseKey, customerName string, expireDate time.Time) (*ent.License, error) {
	newLicense, err := s.client.License.Create().
		SetDomain(domain).
		SetLicenseKey(licenseKey).
		SetCustomerName(customerName).
		SetExpireDate(expireDate).
		SetStatus(1).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return newLicense, nil
}

func (s *LicenseServiceImpl) UpdateLicense(ctx context.Context, id int, domain, licenseKey, customerName string, expireDate time.Time, status int) (*ent.License, error) {
	updateBuilder := s.client.License.UpdateOneID(id)

	if domain != "" {
		updateBuilder.SetDomain(domain)
	}

	if licenseKey != "" {
		updateBuilder.SetLicenseKey(licenseKey)
	}

	if customerName != "" {
		updateBuilder.SetCustomerName(customerName)
	}

	if !expireDate.IsZero() {
		updateBuilder.SetExpireDate(expireDate)
	}

	if status != 0 {
		updateBuilder.SetStatus(status)
	}

	updatedLicense, err := updateBuilder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedLicense, nil
}

func (s *LicenseServiceImpl) DeleteLicense(ctx context.Context, id int) error {
	err := s.client.License.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *LicenseServiceImpl) VerifyLicense(ctx context.Context, domain string) (*ent.License, error) {
	licenseEntity, err := s.client.License.Query().
		Where(license_ent.Domain(domain)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if licenseEntity.ExpireDate.Before(now) {
		_, err = s.client.License.UpdateOneID(licenseEntity.ID).
			SetStatus(2).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	if licenseEntity.Status != 1 {
		return nil, nil
	}

	return licenseEntity, nil
}
