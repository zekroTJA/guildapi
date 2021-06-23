package web

import "github.com/bwmarrin/discordgo"

type guildWrapper struct {
	*discordgo.Guild

	IconURL           string `json:"icon_url"`
	TotalMemberCount  int    `json:"total_member_count"`
	OnlineMemberCount int    `json:"online_member_count"`
}
