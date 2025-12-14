package domain

import (
	"time"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Profile   Profile   `json:"profile"`
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
}

type Profile struct {
	UpdatedAt  time.Time `json:"updated_at"`
	BirthDate  *string   `json:"birth_date,omitempty"`
	LastName   *string   `json:"last_name,omitempty"`
	FirstName  *string   `json:"first_name,omitempty"`
	MiddleName *string   `json:"middle_name,omitempty"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	Gender     Gender    `json:"gender,omitempty"`
}

type Gender struct {
	Value int32  `json:"value"`
	Label string `json:"label"`
}

func genderFromProto(lang string, g userspb.Gender) Gender {
	switch g {
	case userspb.Gender_FEMALE:
		label, _ := i18n.Local.Get(lang, keys.UserGenderFemale)
		return Gender{Value: int32(g), Label: label}
	case userspb.Gender_MALE:
		label, _ := i18n.Local.Get(lang, keys.UserGenderMale)
		return Gender{Value: int32(g), Label: label}
	case userspb.Gender_GENDER_UNSPECIFIED:
		label, _ := i18n.Local.Get(lang, keys.UserGenderUnspecified)
		return Gender{Value: int32(g), Label: label}
	default:
		label, _ := i18n.Local.Get(lang, keys.UserGenderUnspecified)
		return Gender{Value: int32(g), Label: label}
	}
}

func UserFromProto(lang string, u *userspb.User) *User {
	return &User{
		CreatedAt: u.GetCreatedAt().AsTime(),
		UpdatedAt: u.GetUpdatedAt().AsTime(),
		ID:        u.GetId(),
		Email:     u.GetEmail(),
		Phone:     u.Phone,
		Profile: Profile{
			UpdatedAt:  u.Profile.GetUpdatedAt().AsTime(),
			BirthDate:  u.Profile.BirthDate,
			LastName:   u.Profile.LastName,
			FirstName:  u.Profile.FirstName,
			MiddleName: u.Profile.MiddleName,
			AvatarURL:  u.Profile.AvatarUrl,
			Gender:     genderFromProto(lang, u.Profile.GetGender()),
		},
	}
}
