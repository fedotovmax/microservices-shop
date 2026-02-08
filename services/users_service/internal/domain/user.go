package domain

import (
	"time"

	"github.com/fedotovmax/goutils/timeutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserOKResponse struct {
	UID   string
	Email string
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

func (u *User) ToProto() *userspb.User {
	return &userspb.User{
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
		Id:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Profile: &userspb.Profile{
			UpdatedAt:  timestamppb.New(u.Profile.UpdatedAt),
			BirthDate:  timeutils.TimePtrToString(keys.DateFormat, u.Profile.BirthDate),
			LastName:   u.Profile.LastName,
			FirstName:  u.Profile.FirstName,
			MiddleName: u.Profile.MiddleName,
			AvatarUrl:  u.Profile.AvatarURL,
			Gender:     u.Profile.Gender.ToProto(),
		},
	}
}

type GenderValue int8

var (
	GenderUnselected GenderValue = 1
	GenderMale       GenderValue = 2
	GenderFemale     GenderValue = 3
)

func GenderFromProto(fromProto *userspb.GenderValue) *GenderValue {

	if fromProto == nil {
		return nil
	}

	g := GenderValue(*fromProto)

	return &g
}

func (g GenderValue) ToProto() userspb.GenderValue {

	switch g {

	case GenderMale:
		return userspb.GenderValue_GENDER_MALE
	case GenderFemale:
		return userspb.GenderValue_GENDER_FEMALE
	case GenderUnselected:
		return userspb.GenderValue_GENDER_UNSELECTED
	default:
		return userspb.GenderValue_GENDER_UNSPECIFIED
	}

}

func (g GenderValue) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnselected:
		return true
	default:
		return false
	}
}

type Profile struct {
	UpdatedAt  time.Time
	BirthDate  *time.Time
	LastName   *string
	FirstName  *string
	MiddleName *string
	AvatarURL  *string
	Gender     GenderValue
}

type UserPrimaryFields struct {
	ID    string
	Email string
}
