package domain

import (
	"time"

	"github.com/fedotovmax/goutils/timeutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserSessionStatus int8

const (
	UserSessionStatusUnspecified      UserSessionStatus = 0
	UserSessionStatusDeleted          UserSessionStatus = 1
	UserSessionStatusEmailNotVerified UserSessionStatus = 2
	UserSessionStatusBadCredentials   UserSessionStatus = 3
	UserSessionStatusOK               UserSessionStatus = 4
)

func (s UserSessionStatus) ToProto() userspb.UserSessionActionStatus {
	switch s {
	case UserSessionStatusUnspecified:
		return userspb.UserSessionActionStatus_SESSION_STATUS_UNSPECIFIED
	case UserSessionStatusDeleted:
		return userspb.UserSessionActionStatus_SESSION_STATUS_DELETED
	case UserSessionStatusEmailNotVerified:
		return userspb.UserSessionActionStatus_SESSION_STATUS_EMAIL_NOT_VERIFIED
	case UserSessionStatusBadCredentials:
		return userspb.UserSessionActionStatus_SESSION_STATUS_BAD_CREDENTIALS
	case UserSessionStatusOK:
		return userspb.UserSessionActionStatus_SESSION_STATUS_OK
	default:
		return userspb.UserSessionActionStatus_SESSION_STATUS_UNSPECIFIED
	}
}

type UserSessionActionResponse struct {
	Fields UserPrimaryFields
	Status UserSessionStatus
}

func NewUserSessionActionResponse(id, email string, status UserSessionStatus) *UserSessionActionResponse {
	return &UserSessionActionResponse{
		Status: status,
		Fields: UserPrimaryFields{
			ID:    id,
			Email: email,
		},
	}
}

func (sr *UserSessionActionResponse) ToProto() *userspb.UserSessionActionResponse {

	var email *string
	var userId *string

	if sr.Fields.Email != "" {
		email = &sr.Fields.Email
	}

	if sr.Fields.ID != "" {
		userId = &sr.Fields.ID
	}

	return &userspb.UserSessionActionResponse{
		UserSessionActionStatus: sr.Status.ToProto(),
		Email:                   email,
		UserId:                  userId,
	}
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

type GenderValue int8

var (
	GenderInvalid    GenderValue = -1
	GenderUnselected GenderValue = 1
	GenderMale       GenderValue = 2
	GenderFemale     GenderValue = 3
)

func GenderFromProto(g *userspb.GenderValue) *GenderValue {

	if g == nil {
		return nil
	}

	switch *g {
	case userspb.GenderValue_GENDER_MALE:
		male := GenderMale
		return &male
	case userspb.GenderValue_GENDER_FEMALE:
		female := GenderFemale
		return &female
	case userspb.GenderValue_GENDER_UNSELECTED:
		unselected := GenderUnselected
		return &unselected
	default:
		invalid := GenderInvalid
		return &invalid
	}
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

type EmailVerifyLink struct {
	Link          string
	UserID        string
	LinkExpiresAt time.Time
}

type UserPrimaryFields struct {
	ID    string
	Email string
}
