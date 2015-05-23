package game_test

import (
	"github.com/nkcraddock/numzero/game"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/redis.v3"
)

var _ = Describe("game.redisStore integration tests", func() {
	var store *game.RedisStore
	chad := &game.Player{Name: "Chad"}
	roger := &game.Player{Name: "Roger"}
	got_powerup := &game.Rule{"powerup", "got the powerup", 5}
	won_thegame := &game.Rule{"wonthegame", "won the game", 20}

	BeforeEach(func() {
		options := &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       10,
		}

		store = game.NewRedisStore(options)
		store.FlushDb()
	})

	Context("Rules", func() {
		It("saves a rule", func() {
			err := store.SaveRule(got_powerup)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("retrieves a rule", func() {
			store.SaveRule(got_powerup)

			r, err := store.GetRule("powerup")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(r).ShouldNot(BeNil())
			Ω(r.Points).Should(Equal(5))
		})

		It("retrieves a list of rules", func() {
			store.SaveRule(got_powerup)
			store.SaveRule(won_thegame)

			r, err := store.ListRules()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(r).ShouldNot(BeNil())
			Ω(r).Should(HaveLen(2))
		})
	})

	Context("Players", func() {
		It("saves a player", func() {
			err := store.SavePlayer(chad)
			Ω(err).ShouldNot(HaveOccurred())

			p, err := store.GetPlayer("chad")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(p).ShouldNot(BeNil())
			Ω(p.Name).Should(Equal("Chad"))
		})

		It("retrieves a player", func() {
			store.SavePlayer(chad)

			p, err := store.GetPlayer("chad")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(p).ShouldNot(BeNil())
		})

		It("returns an error if player doesnt exist", func() {
			_, err := store.GetPlayer("chad")
			Ω(err).Should(Equal(game.ErrorNotFound))
		})

		It("retrieves a list of players", func() {
			store.SavePlayer(chad)
			store.SavePlayer(roger)

			players, err := store.ListPlayers()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(players).Should(HaveLen(2))
		})

	})
})