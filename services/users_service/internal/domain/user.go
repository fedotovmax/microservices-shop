package domain

import (
	"time"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *User) ToProto(locale string) *userspb.User {

	return &userspb.User{
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
		Id:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Profile: &userspb.Profile{
			UpdatedAt:  timestamppb.New(u.Profile.UpdatedAt),
			BirthDate:  utils.TimePtrToString(keys.DateFormat, u.Profile.BirthDate),
			LastName:   u.Profile.LastName,
			FirstName:  u.Profile.FirstName,
			MiddleName: u.Profile.MiddleName,
			AvatarUrl:  u.Profile.AvatarURL,
			Gender:     u.Profile.Gender.ToProto(locale),
		},
	}
}

type User struct {
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
	Profile         Profile
	Phone           *string
	ID              string
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	IsPhoneVerified bool
}

type Gender string

func (g *Gender) String() string {

	if g == nil {
		return ""
	}

	return string(*g)

}

func GenderFromProto(g *string) *Gender {
	if g == nil {
		return nil
	}

	ng := Gender(*g)

	return &ng
}

func (g Gender) ToProto(locale string) string {
	switch g {
	case Male:
		label, _ := i18n.Local.Get(locale, keys.UserGenderMale)
		return label
	case Female:
		label, _ := i18n.Local.Get(locale, keys.UserGenderFemale)
		return label
	case Unspecified:
		label, _ := i18n.Local.Get(locale, keys.UserGenderUnspecified)
		return label
	default:
		label, _ := i18n.Local.Get(locale, keys.UserGenderUnspecified)
		return label
	}
}

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female, Unspecified:
		return true
	default:
		return false
	}
}

var (
	Male          Gender = "male"
	Female        Gender = "female"
	Unspecified   Gender = "unspecified"
	InvalidGender Gender = ""
)

type Profile struct {
	UpdatedAt  time.Time
	BirthDate  *time.Time
	LastName   *string
	FirstName  *string
	MiddleName *string
	AvatarURL  *string
	Gender     Gender
}

type EmailVerifyLink struct {
	Link           string
	UserID         string
	ValidityPeriod time.Time
}
