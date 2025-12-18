package domain

import (
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
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
	Gender     string    `json:"gender,omitempty"`
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
			Gender:     u.Profile.GetGender(),
		},
	}
}
