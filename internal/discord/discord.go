package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/dgrs"
)

type Session struct {
	session *discordgo.Session
	state   *dgrs.State
}

func New(token string, redisClient *redis.Client) (s *Session, err error) {
	s = new(Session)
	if s.session, err = discordgo.New("Bot " + token); err != nil {
		return
	}
	s.session.StateEnabled = false
	s.session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildBans |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildPresences

	s.session.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.Ready) {
		for _, g := range e.Guilds {
			ms, err := s.state.Members(g.ID, true)
			if err != nil {
				logrus.WithError(err).WithField("gid", g.ID).Error("Failed fetching guild members")
			} else {
				logrus.WithField("gid", g.ID).WithField("n", len(ms)).Info("Fetched guild members")
			}
		}
	})

	s.state, err = dgrs.New(dgrs.Options{
		RedisClient:    redisClient,
		DiscordSession: s.session,
		FetchAndStore:  true,
		FlushOnStartup: true,
		KeyPrefix:      "guildapi",
	})
	return
}

func (s *Session) Open() error {
	return s.session.Open()
}

func (s *Session) Close() error {
	return s.session.Close()
}

func (s *Session) Guild(id string) (g *discordgo.Guild, err error) {
	return s.state.Guild(id)
}
