package models

type Relationship struct {
	UserId             string `bson:"user_id" json:"user_id"`
	UserRelationshipId string `bson:"user_relationship_id" json:"user_relationship_id"`
}

type ReponseGetRelationship struct {
	Status bool `json:"status"`
}
