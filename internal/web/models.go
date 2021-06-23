package web

import "github.com/bwmarrin/discordgo"

type guildWrapper struct {
	*discordgo.Guild

	TotalMemberCount  int `json:"total_member_count"`
	OnlineMemberCount int `json:"online_member_count"`
}
