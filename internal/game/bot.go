package game

import (
	"math/rand"
	"time"

	"github.com/OpenBachelor/OpenBachelorSS/internal/session"
	"github.com/OpenBachelor/OpenBachelorSS/pkg/contract"
)

// BotPlayer represents an AI-controlled player that participates
// in the game without a real TCP session.
type BotPlayer struct {
	Status EnemyDuelGamePlayerStatus
}

var botNames = []string{
	"博士", "阿米娅", "凯尔希", "能天使", "德克萨斯",
	"拉普兰德", "银灰", "陈", "星熊", "煌",
	"塞雷娅", "伊芙利特", "艾雅法拉", "推进之王", "闪灵",
	"夜莺", "安洁莉娜", "斯卡蒂", "黑", "麦哲伦",
	"莫斯提马", "煌", "阿", "年", "刻俄柏",
}

var botAvatars = []string{
	"avatar_def_01", "avatar_def_02", "avatar_def_03", "avatar_def_04",
	"char_002_amiya", "char_103_angel", "char_102_texas", "char_144_silverash",
	"char_010_chen", "char_147_shining", "char_179_cgbird", "char_291_aglin",
	"char_263_skadi", "char_340_halfm", "char_248_mgllan", "char_213_mostma",
}

func generateBotInfo() (string, string) {
	name := botNames[rand.Intn(len(botNames))]
	avatar := botAvatars[rand.Intn(len(botAvatars))]
	return name, avatar
}

// createBots fills remaining player slots with bot players.
// Should be called when the game transitions out of WaitingState.
func (gm *EnemyDuelGame) createBots() {
	maxPlayers := gm.getMaxNumPlayer()
	numBots := maxPlayers - gm.nextInternalPlayerID
	if numBots <= 0 {
		return
	}

	gm.bots = make([]*BotPlayer, numBots)
	for i := 0; i < numBots; i++ {
		bot := &BotPlayer{}
		bot.Status.internalPlayerID = gm.nextInternalPlayerID
		gm.nextInternalPlayerID++

		gm.initPlayerStatusForBot(bot)
		gm.bots[i] = bot
	}
}

func (gm *EnemyDuelGame) initPlayerStatusForBot(bot *BotPlayer) {
	bot.Status.PlayerID = bot.Status.getExternalPlayerID()
	
	name, avatar := generateBotInfo()
	bot.Status.NickName = name
	bot.Status.AvatarID = avatar

	switch gm.ModeID {
	case "multiOperationMatch":
		bot.Status.Money = 10000
	default:
		bot.Status.Money = 1
		bot.Status.ShieldState = 2
	}
	bot.Status.IsReady = true
}

// makeBetDecision generates a random betting decision.
// Side: 0b01 (LEFT) or 0b10 (RIGHT), with a small chance of AllIn.
func (bot *BotPlayer) makeBetDecision() {
	// 0b01 = LEFT, 0b10 = RIGHT
	bot.Status.Side = uint8(1 << rand.Intn(2))

	if rand.Float64() < 0.1 {
		bot.Status.AllIn = 1
	} else {
		bot.Status.AllIn = 0
	}
}

// broadcastBotBets has all bots place bets and broadcasts to real sessions.
func (gm *EnemyDuelGame) broadcastBotBets(sessions map[*session.Session]*EnemyDuelSessionGameStatus, forceExitTime time.Time) {
	for _, bot := range gm.bots {
		bot.makeBetDecision()
		for session := range sessions {
			session.SendMessage(
				contract.NewS2CEnemyDuelClientStateMessageForBet(
					2, gm.round, forceExitTime,
					bot.Status.getExternalPlayerID(),
					bot.Status.Side,
					bot.Status.AllIn,
					bot.Status.Streak,
				),
			)
		}
	}
}

// settleAllBots makes all bots report the given side as their settle result.
// Called when a real player sends RoundSettle to keep bots in sync.
func (gm *EnemyDuelGame) settleAllBots(side uint8) {
	for _, bot := range gm.bots {
		bot.Status.ReportSide = side
	}
}
