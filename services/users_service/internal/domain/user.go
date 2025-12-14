package domain

import (
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *User) ToProto() *userspb.User {

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
			Gender:     u.Profile.Gender.ToProto(),
		},
	}
}

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Profile      Profile
	ID           string
	Email        string
	Phone        *string
	PasswordHash string
}

type Gender string

func (g *Gender) String() string {

	if g == nil {
		return ""
	}

	return string(*g)

}

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female:
		return true
	default:
		return false
	}
}

func GenderFromProto(protoG *userspb.Gender) *Gender {

	if protoG == nil {
		return nil
	}

	switch *protoG {
	case userspb.Gender_MALE:
		g := Male
		return &g
	case userspb.Gender_FEMALE:
		g := Female
		return &g
	case userspb.Gender_GENDER_UNSPECIFIED:
		return nil
	default:
		g := InvalidGender
		return &g
	}

}

func (g *Gender) ToProto() *userspb.Gender {

	if g == nil {
		return userspb.Gender_GENDER_UNSPECIFIED.Enum()
	}

	switch *g {
	case Male:
		return userspb.Gender_MALE.Enum()
	case Female:
		return userspb.Gender_FEMALE.Enum()
	default:
		return userspb.Gender_GENDER_UNSPECIFIED.Enum()
	}

}

var (
	Male          Gender = "male"
	Female        Gender = "female"
	InvalidGender Gender = ""
)

type Profile struct {
	UpdatedAt  time.Time
	BirthDate  *time.Time
	LastName   *string
	FirstName  *string
	MiddleName *string
	AvatarURL  *string
	Gender     *Gender
}
