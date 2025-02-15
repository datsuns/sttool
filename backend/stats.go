package backend

import (
	"fmt"
	"io"
	"time"
)

type UserName string
type ChannelPointTitle string

type PeriodStats struct {
	Started  time.Time
	Finished time.Time
}

type FollowStats struct {
	Users []UserName
}

type ViewerStats struct {
	Time  time.Time
	Total int
}

type ChatEntry struct {
	Time time.Time
	User UserName
	Text string
}

type ChatStats struct {
	Total   int
	History []ChatEntry
}

type BitsRecord struct {
	Bits  int
	Times int
}

type CheerStats struct {
	TotalBits int
	History   map[UserName]BitsRecord
}

type SubGiftStats struct {
	TotalGifts int
	History    map[UserName]int
}

type SubGiftReceived struct {
	History map[UserName]int
}

type SubScriptionEntry struct {
	Tier string
}

type SubScriptionStats struct {
	Entry map[UserName]SubScriptionEntry
}

type ChannelPointStats struct {
	TotalTimes int
	Record     map[UserName]int
}

type RaidEntry struct {
	From    UserName
	Viewers int
}

type RaidStats struct {
	History []RaidEntry
}

type GigantifiedEmoteHistory struct {
	Times   int
	History map[UserName]int
}

type MessegeEffectHistory struct {
	Times   int
	History map[UserName]int
}

type PowerUpStats struct {
	GigantifiedEmoteHistory GigantifiedEmoteHistory
	MessegeEffectHistory    MessegeEffectHistory
}

type TwitchStats struct {
	InStreaming       bool
	LastPeriod        PeriodStats
	FollowStats       FollowStats
	ChatStats         ChatStats
	CheerStats        CheerStats
	SubScriptionStats SubScriptionStats
	SubGiftStats      SubGiftStats
	SubGiftReceived   SubGiftReceived
	ViewersHistory    []ViewerStats
	ChannelPoinsts    ChannelPointStats
	RaidStats         RaidStats
	PowerUpStats      PowerUpStats
}

func NewTwitchStats() *TwitchStats {
	ret := &TwitchStats{}
	ret.Clear()
	return ret
}

func (t *TwitchStats) Clear() {
	t.InStreaming = false
	t.FollowStats.Users = []UserName{}
	t.ChatStats = ChatStats{
		Total: 0,
	}
	t.CheerStats = CheerStats{
		TotalBits: 0,
		History:   map[UserName]BitsRecord{},
	}
	t.SubScriptionStats = SubScriptionStats{
		Entry: map[UserName]SubScriptionEntry{},
	}
	t.SubGiftStats = SubGiftStats{
		TotalGifts: 0,
		History:    map[UserName]int{},
	}
	t.SubGiftReceived = SubGiftReceived{
		History: map[UserName]int{},
	}
	t.ViewersHistory = []ViewerStats{}
	t.ChannelPoinsts = ChannelPointStats{
		TotalTimes: 0,
		Record:     map[UserName]int{},
	}
	t.RaidStats = RaidStats{
		History: []RaidEntry{},
	}
	t.PowerUpStats = PowerUpStats{
		GigantifiedEmoteHistory: GigantifiedEmoteHistory{
			Times:   0,
			History: map[UserName]int{},
		},
		MessegeEffectHistory: MessegeEffectHistory{
			Times:   0,
			History: map[UserName]int{},
		},
	}
}

func (t *TwitchStats) String(topIndent, namePrefix string) string {
	raidTimes, _ := t.LoadRaidResult()
	started := t.LastPeriod.Started.Format("2006/01/02 15:04:05")
	finished := t.LastPeriod.Finished.Format("2006/01/02 15:04:05")
	followResult := fmt.Sprintf("%v新規フォロー: %v人\n", topIndent, len(t.FollowStats.Users))
	for _, u := range t.FollowStats.Users {
		followResult += fmt.Sprintf("%v  %v%vさん\n", topIndent, namePrefix, u)
	}
	chanepoResult := fmt.Sprintf("%vチャネポ総回数: %v\n", topIndent, t.LoadChannelPointTotal())
	for name, times := range t.LoadChannelPointHistory() {
		chanepoResult += fmt.Sprintf("%v  %v%vさん: %v回\n", topIndent, namePrefix, name, times)
	}
	subscResult := fmt.Sprintf("%v新規サブスク: %v人\n", topIndent, len(t.LoadSubScribed()))
	for name := range t.LoadSubscriptonHistory() {
		subscResult += fmt.Sprintf("%v  %v%vさん\n", topIndent, namePrefix, name)
	}
	subGifResult := fmt.Sprintf("%v総サブギフ個数: %v個\n", topIndent, t.LoadSubGiftTotal())
	for name, times := range t.LoadSubGiftHistory() {
		subGifResult += fmt.Sprintf("%v  %v%vさん(%v個)\n", topIndent, namePrefix, name, times)
	}
	subGifRecvResult := fmt.Sprintf("%v  >> サブギフ受け取った: %v人\n", topIndent, len(t.LoadSubGifted()))
	for name := range t.LoadSubGifted() {
		subGifRecvResult += fmt.Sprintf("%v    %v%vさん\n", topIndent, namePrefix, name)
	}
	cheerResult := fmt.Sprintf("%vビッツ: %v\n", topIndent, t.LoadCheerTotal())
	for name, bitsRecord := range t.LoadCheerHistory() {
		cheerResult += fmt.Sprintf("%v  %v%vさん(%v ビッツ)\n", topIndent, namePrefix, name, bitsRecord.Bits)
	}
	raidResult := fmt.Sprintf("%vレイド: %v回\n", topIndent, raidTimes)
	for _, e := range t.LoadRaidHistory() {
		raidResult += fmt.Sprintf("%v  %v%vさん\n", topIndent, namePrefix, e.From)
	}
	gigantifiedEmoteResult := fmt.Sprintf("%v巨大化スタンプ: %v回\n", topIndent, t.LoadGigantifiedEmoteTimes())
	for k, v := range t.LoadGigantifiedEmoteHistory() {
		gigantifiedEmoteResult += fmt.Sprintf("%v  %v%vさん : %v回\n", topIndent, namePrefix, k, v)
	}
	messageEffectResult := fmt.Sprintf("%vメッセージエフェクト: %v回\n", topIndent, t.LoadMessageEffectTimes())
	for k, v := range t.LoadMessageEffectHistory() {
		messageEffectResult += fmt.Sprintf("%v  %v%vさん : %v回\n", topIndent, namePrefix, k, v)
	}
	return fmt.Sprintf(
		"------------------------------------------------------------\n"+
			"%v配信時間: %v ~ %v\n"+
			"%v"+
			"%v"+
			"%v"+
			"%v"+
			"%v"+
			"%v"+
			"%v"+
			"%v"+
			"%v",
		topIndent, started, finished,
		followResult,
		chanepoResult,
		subscResult,
		subGifResult,
		subGifRecvResult,
		cheerResult,
		raidResult,
		gigantifiedEmoteResult,
		messageEffectResult,
	)
}

func (t *TwitchStats) Dump(w io.Writer, topIndent, namePrefix string) {
	w.Write([]byte(t.String(topIndent, namePrefix) + "\n"))
}

func (t *TwitchStats) StreamStarted() {
	t.Clear()
	t.InStreaming = true
	t.LastPeriod.Started = time.Now()
}

func (t *TwitchStats) StreamFinished() {
	t.LastPeriod.Finished = time.Now()
	t.InStreaming = false
}

func (t *TwitchStats) Follow(user UserName) {
	if t.InStreaming == false {
		return
	}
	t.FollowStats.Users = append(t.FollowStats.Users, user)
}

func (t *TwitchStats) Chat(user UserName, text string) {
	if t.InStreaming == false {
		return
	}
	t.ChatStats.Total += 1
	t.ChatStats.History = append(t.ChatStats.History, ChatEntry{Time: time.Now(), User: user, Text: text})
}

func (t *TwitchStats) ChannelPoint(user UserName, title ChannelPointTitle) {
	if t.InStreaming == false {
		return
	}
	t.ChannelPoinsts.TotalTimes += 1
	if _, exists := t.ChannelPoinsts.Record[user]; exists {
		t.ChannelPoinsts.Record[user] += 1
	} else {
		t.ChannelPoinsts.Record[user] = 1
	}
}

func (t *TwitchStats) Cheer(user UserName, n int) {
	t.CheerStats.TotalBits += n
	if v, exists := t.CheerStats.History[user]; exists {
		v.Bits += n
		v.Times += 1
		t.CheerStats.History[user] = v
	} else {
		t.CheerStats.History[user] = BitsRecord{Bits: n, Times: 1}
	}
}

func (t *TwitchStats) SubGift(user UserName, n int) {
	t.SubGiftStats.TotalGifts += n
	if v, exists := t.SubGiftStats.History[user]; exists {
		v += n
		t.SubGiftStats.History[user] = v
	} else {
		t.SubGiftStats.History[user] = n
	}
}

func (t *TwitchStats) SubGifted(user UserName, tier string) {
	if v, exists := t.SubGiftReceived.History[user]; exists {
		v += 1
		t.SubGiftReceived.History[user] = v
	} else {
		t.SubGiftReceived.History[user] = 1
	}
}

func (t *TwitchStats) SubScribe(user UserName, tier string) {
	if v, exists := t.SubScriptionStats.Entry[user]; exists {
		v.Tier = tier
		t.SubScriptionStats.Entry[user] = v
	} else {
		t.SubScriptionStats.Entry[user] = SubScriptionEntry{Tier: tier}
	}
}

func (t *TwitchStats) Raid(from UserName, viewers int) {
	t.RaidStats.History = append(
		t.RaidStats.History,
		RaidEntry{From: from, Viewers: viewers},
	)
}

func (t *TwitchStats) GigantifiedEmote(from UserName) {
	t.PowerUpStats.GigantifiedEmoteHistory.Times += 1
	if _, exists := t.PowerUpStats.GigantifiedEmoteHistory.History[from]; exists {
		t.PowerUpStats.GigantifiedEmoteHistory.History[from] += 1
	} else {
		t.PowerUpStats.GigantifiedEmoteHistory.History[from] = 1
	}
}

func (t *TwitchStats) MessageEffect(from UserName) {
	t.PowerUpStats.MessegeEffectHistory.Times += 1
	if _, exists := t.PowerUpStats.MessegeEffectHistory.History[from]; exists {
		t.PowerUpStats.MessegeEffectHistory.History[from] += 1
	} else {
		t.PowerUpStats.MessegeEffectHistory.History[from] = 1
	}
}

// --- loader

func (t *TwitchStats) LoadPeriod() time.Duration {
	return t.LastPeriod.Finished.Sub(t.LastPeriod.Started)
}

func (t *TwitchStats) LoadNChats() int {
	return t.ChatStats.Total
}

func (t *TwitchStats) LoadChatHistory() []ChatEntry {
	return t.ChatStats.History
}

func (t *TwitchStats) LoadCheerTotal() int {
	return t.CheerStats.TotalBits
}

func (t *TwitchStats) LoadCheerHistory() map[UserName]BitsRecord {
	return t.CheerStats.History
}

func (t *TwitchStats) LoadSubGiftTotal() int {
	return t.SubGiftStats.TotalGifts
}

func (t *TwitchStats) LoadSubGiftHistory() map[UserName]int {
	return t.SubGiftStats.History
}

func (t *TwitchStats) LoadSubGifted() map[UserName]int {
	return t.SubGiftReceived.History
}

func (t *TwitchStats) LoadSubScribed() map[UserName]SubScriptionEntry {
	return t.SubScriptionStats.Entry
}

func (t *TwitchStats) LoadChannelPointTotal() int {
	return t.ChannelPoinsts.TotalTimes
}

func (t *TwitchStats) LoadChannelPointHistory() map[UserName]int {
	return t.ChannelPoinsts.Record
}

func (t *TwitchStats) LoadChannelPointTimes(user UserName) int {
	if _, exists := t.ChannelPoinsts.Record[user]; exists {
		return t.ChannelPoinsts.Record[user]
	} else {
		return 0
	}
}

func (t *TwitchStats) LoadSubscriptonHistory() map[UserName]SubScriptionEntry {
	return t.SubScriptionStats.Entry
}

func (t *TwitchStats) LoadRaidResult() (int, int) {
	times := len(t.RaidStats.History)
	total := 0
	for _, e := range t.RaidStats.History {
		total += e.Viewers
	}
	return times, total
}

func (t *TwitchStats) LoadRaidHistory() []RaidEntry {
	return t.RaidStats.History
}

func (t *TwitchStats) LoadGigantifiedEmoteTimes() int {
	return t.PowerUpStats.GigantifiedEmoteHistory.Times
}

func (t *TwitchStats) LoadGigantifiedEmoteHistory() map[UserName]int {
	return t.PowerUpStats.GigantifiedEmoteHistory.History
}

func (t *TwitchStats) LoadMessageEffectTimes() int {
	return t.PowerUpStats.MessegeEffectHistory.Times
}

func (t *TwitchStats) LoadMessageEffectHistory() map[UserName]int {
	return t.PowerUpStats.MessegeEffectHistory.History
}
