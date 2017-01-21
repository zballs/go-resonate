package api

import . "github.com/zballs/go_resonate/util"

type UserInfo struct {
	Privkey *PrivateKey `json:"private_key"`
	Pubkey  *PublicKey  `json:"public_key"`
	UserId  string      `json:"user_id"`
}

func NewUserInfo(userId string, priv *PrivateKey, pub *PublicKey) *UserInfo {
	return &UserInfo{
		UserId:  userId,
		Privkey: priv,
		Pubkey:  pub,
	}
}

type ProjectInfo struct {
	ProjectId string   `json:"project_id"`
	SongIds   []string `json:"song_ids"`
}

func NewProjectInfo(projectId string, songIds []string) *ProjectInfo {
	return &ProjectInfo{
		ProjectId: projectId,
		SongIds:   songIds,
	}
}
