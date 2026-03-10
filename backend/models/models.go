package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text          string             `bson:"text" json:"text"`
	ImageURL      string             `bson:"image_url" json:"image_url"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer int                `bson:"correct_answer" json:"correct_answer"`
	TimeLimit     int                `bson:"time_limit" json:"time_limit"` // in seconds
	Type          string             `bson:"type" json:"type"`
	Points        string             `bson:"points" json:"points"`
	AnswerType    string             `bson:"answer_type" json:"answer_type"`
}

type Game struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Questions []Question         `bson:"questions" json:"questions"`
	RoomID    string             `bson:"room_id" json:"room_id"`
	Active    bool               `bson:"active" json:"active"`
}

type Player struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"` // Hashed
}

type QuizHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuizID       primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	QuizTitle    string             `bson:"quiz_title" json:"quiz_title"`
	RoomID       string             `bson:"room_id" json:"room_id"`
	PlayedAt     primitive.DateTime `bson:"played_at" json:"played_at"`
	PlayersCount int                `bson:"players_count" json:"players_count"`
	Leaderboard  []PlayerResult     `bson:"leaderboard" json:"leaderboard"`
}

type PlayerResult struct {
	Username string `bson:"username" json:"username"`
	Score    int    `bson:"score" json:"score"`
}
