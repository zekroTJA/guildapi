package web

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zekrotja/guildapi/internal/discord"
)

type Web struct {
	session *discord.Session
	app     *fiber.App
}

func New(session *discord.Session, alloworigins string) (w *Web) {
	w = new(Web)

	w.session = session

	w.app = fiber.New(fiber.Config{
		GETOnly: true,
	})

	w.app.Use(cors.New(cors.Config{
		AllowHeaders: "*",
		AllowMethods: "OPTIONS, GET",
		AllowOrigins: alloworigins,
	}))

	w.app.Get("/guilds/:id", w.getGuild)

	return
}

func (w *Web) Open(address string) error {
	return w.app.Listen(address)
}

func (w *Web) getGuild(c *fiber.Ctx) error {
	id := c.Params("id")

	_g, err := w.session.Guild(id, true)
	if err != nil {
		return err
	}

	g := &guildWrapper{Guild: _g}

	if g.Members != nil {
		for _, m := range g.Members {
			if !m.User.Bot {
				g.TotalMemberCount++
			}
		}
	}

	presences, err := w.session.Presences(g.ID)
	if err != nil {
		return err
	}
	for _, p := range presences {
		if !p.User.Bot && p.Status != discordgo.StatusOffline {
			g.OnlineMemberCount++
		}
	}

	include := strings.Split(c.Query("include", ""), ",")

	if !contains(include, "all") {
		if !contains(include, "roles") {
			g.Roles = make([]*discordgo.Role, 0)
		}
		if !contains(include, "emojis") {
			g.Emojis = make([]*discordgo.Emoji, 0)
		}
		if !contains(include, "members") {
			g.Members = make([]*discordgo.Member, 0)
		}
		if !contains(include, "channels") {
			g.Channels = make([]*discordgo.Channel, 0)
		}
	}

	return c.JSON(g)
}

func contains(s []string, v string) bool {
	v = trimLower(v)
	for _, st := range s {
		if v == trimLower(st) {
			return true
		}
	}
	return false
}

func trimLower(v string) string {
	return strings.ToLower(strings.Trim(v, " \t"))
}
