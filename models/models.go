package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=30" bson:"name"`
	Password      *string            `json:"password" validate:"required,min=6" bson:"password"`
	Email         *string            `json:"email" validate:"email,required" bson:"email"`
	Role          int                `json:"role" bson:"role"`
	Avatar        *string            `json:"avatar" bson:"avatar"`
	AvatarPath    *string            `json:"avatar_path" bson:"avatar_path"`
	Token         *string            `json:"token" bson:"token"`
	Refresh_Token *string            `json:"refresh_token" bson:"refresh_token"`
	User_ID       string             `json:"user_id" bson:"user_id"`
	Created_At    time.Time          `json:"created_at" bson:"created_at"`
	Updated_At    time.Time          `json:"updated_at" bson:"updated_at"`
}

type Service struct {
	Service_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Title       *string            `json:"title" bson:"title"`
	Description *string            `json:"description" bson:"description"`
	Image       *string            `json:"image" bson:"image"`
	ImagePath   *string            `json:"image_path" bson:"image_path"`
	T1          *string            `json:"t1" bson:"t1"`
	T2          *string            `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Images struct {
	Image     *string `json:"image" bson:"image"`
	ImagePath *string `json:"image_path" bson:"image_path"`
}

type Project struct {
	Project_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Title       *string            `json:"title" bson:"title"`
	Description *string            `json:"description" bson:"description"`
	DemoLink    *string            `json:"demo_link" bson:"demo_link"`
	Tech        *string            `json:"tech" bson:"tech"`
	Images      []*Images          `json:"images" bson:"images"`
	T1          *string            `json:"t1" bson:"t1"`
	T2          *string            `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Remark struct {
	Remark_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Name       *string            `json:"name" bson:"name"`
	Role       *string            `json:"role" bson:"role"`
	Image      *string            `json:"image" bson:"image"`
	ImagePath  *string            `json:"image_path" bson:"image_path"`
	Remark     *string            `json:"remark" bson:"remark"`
	T1         *string            `json:"t1" bson:"t1"`
	T2         *string            `json:"t2" bson:"t2"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	Updated_At time.Time          `json:"updated_at" bson:"updated_at"`
}

type Message struct {
	Message_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Name        *string            `json:"name" bson:"name"`
	Email       *string            `json:"email" bson:"email"      validate:"email,required"`
	Phone       *string            `json:"phone" bson:"phone"`
	CompanyName *string            `json:"company_name" bson:"company_name"`
	Message     *string            `json:"message" bson:"message"`
	T1          *string            `json:"t1" bson:"t1"`
	T2          *string            `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}
