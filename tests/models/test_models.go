package models

import (
	"time"

	"github.com/crashana/go-eloquent"
)

// UserModel - Test model for users
type UserModel struct {
	*eloquent.BaseModel

	// Struct fields for direct access
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"password" db:"password"`
	EmailVerifiedAt time.Time `json:"email_verified_at" db:"email_verified_at"`
	IsAdmin         bool      `json:"is_admin" db:"is_admin"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at" db:"deleted_at"`
}

// NewUser creates a new UserModel instance
func NewUser() *UserModel {
	user := &UserModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	user.Table("users").
		PrimaryKey("id").
		Fillable("name", "email", "password", "is_admin", "status").
		Hidden("password", "remember_token").
		Casts(map[string]string{
			"email_verified_at": "datetime",
			"is_admin":          "bool",
			"created_at":        "datetime",
			"updated_at":        "datetime",
			"deleted_at":        "datetime",
		})

	// Set the parent model reference for attribute syncing
	user.SetParentModel(user)

	return user
}

// Global static instance for User model
var User = eloquent.NewModelStatic(func() *UserModel {
	return NewUser()
})

// PostModel - Test model for posts
type PostModel struct {
	*eloquent.BaseModel

	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	UserID    string    `json:"user_id" db:"user_id"`
	Published bool      `json:"published" db:"published"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewPost creates a new PostModel instance
func NewPost() *PostModel {
	post := &PostModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	post.Table("posts").
		PrimaryKey("id").
		Fillable("title", "content", "user_id", "published").
		Casts(map[string]string{
			"published":  "bool",
			"created_at": "datetime",
			"updated_at": "datetime",
		})

	// Set the parent model reference for attribute syncing
	post.SetParentModel(post)

	return post
}

// Define relationships for PostModel
func (p *PostModel) Author() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(p)
	return rb.BelongsTo("author", "UserModel")
}

// Global static instance for Post model
var Post = eloquent.NewModelStatic(func() *PostModel {
	return NewPost()
})

// ProfileModel - Test model for user profiles
type ProfileModel struct {
	*eloquent.BaseModel

	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Bio       string    `json:"bio" db:"bio"`
	Avatar    string    `json:"avatar" db:"avatar"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewProfile creates a new ProfileModel instance
func NewProfile() *ProfileModel {
	profile := &ProfileModel{
		BaseModel: eloquent.NewBaseModel(),
	}

	profile.Table("profiles").
		PrimaryKey("id").
		Fillable("user_id", "bio", "avatar").
		Casts(map[string]string{
			"created_at": "datetime",
			"updated_at": "datetime",
		})

	// Set the parent model reference for attribute syncing
	profile.SetParentModel(profile)

	return profile
}

// Define relationships for ProfileModel
func (p *ProfileModel) User() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(p)
	return rb.BelongsTo("user", "UserModel")
}

// Global static instance for Profile model
var Profile = eloquent.NewModelStatic(func() *ProfileModel {
	return NewProfile()
})

// Define relationships for UserModel
func (u *UserModel) Posts() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(u)
	return rb.HasMany("posts", "PostModel")
}

func (u *UserModel) Profile() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(u)
	return rb.HasOne("profile", "ProfileModel")
}
