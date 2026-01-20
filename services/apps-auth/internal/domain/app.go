package domain

import (
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/appsauthpb"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain/errs"
)

type ApplicationType uint8

const (
	ApplicationTypeInvalid  ApplicationType = 0
	ApplicationTypeCustomer ApplicationType = 1
	ApplicationTypePartner  ApplicationType = 2
	ApplicationTypeAdmin    ApplicationType = 3
)

func ApplicationTypeFromProto(t appsauthpb.ApplicationType) ApplicationType {
	switch t {
	case appsauthpb.ApplicationType_APPLICATION_TYPE_ADMIN:
		return ApplicationTypeAdmin
	case appsauthpb.ApplicationType_APPLICATION_TYPE_CUSTOMER:
		return ApplicationTypeCustomer
	case appsauthpb.ApplicationType_APPLICATION_TYPE_PARTNER:
		return ApplicationTypePartner
	default:
		return ApplicationTypeInvalid
	}
}

func (t ApplicationType) IsValid() error {
	switch t {
	case ApplicationTypeCustomer, ApplicationTypePartner, ApplicationTypeAdmin:
		return nil
	default:
		return errs.ErrInvalidAppType
	}
}

type App struct {
	CreatedAt time.Time       `json:"created_at"`
	Name      string          `json:"name"`
	Type      ApplicationType `json:"type"`
}
