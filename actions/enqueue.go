package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var remindersQueue = make(chan *Reminder, 20)

type Reminder struct {
	Time      time.Time
	Message   string
	ChannelID string
}

func EnqueueReminder(reminder *Reminder) error {
	fmt.Println("Enfileirando lembrete para", reminder.Time)
	fmt.Println("A fila tem", len(remindersQueue), "lembretes agendados")

	select {
	case remindersQueue <- reminder:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func ProcessReminders(s *discordgo.Session) {
	for {
		reminder := <-remindersQueue
		fmt.Println("Lembrete agendado para", reminder.Time)

		duration := time.Until(reminder.Time)
		if duration > 0 {
			<-time.After(duration)
		}

		_, err := s.ChannelMessageSend(reminder.ChannelID, reminder.Message)
		if err != nil {
			fmt.Println("Erro ao enviar a mensagem:", err)
		}
	}
}

func EnviarMensagemNoChat1(s *discordgo.Session, m *discordgo.MessageCreate, trecho string) {
	if m.Author.Bot {
		EnviarMensagem(s, m, "Eu não respondo a bots! HÁ HÁ HÁ")
		return
	}

	voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil || voiceState == nil {
		EnviarMensagem(s, m, "Você precisa estar em um canal de voz para usar esse comando :(")
		return
	}

	args := strings.Split(trecho, " ")
	if len(args) < 3 {
		return
	}
	timeStr := args[2]
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		fmt.Println("Erro ao converter a hora:", err)
		return
	}

	reminder := &Reminder{
		Time:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), t.Hour(), t.Minute(), 0, 0, time.Now().Location()),
		Message:   "Hora da daily!!",
		ChannelID: voiceState.ChannelID,
	}

	if err := EnqueueReminder(reminder); err != nil {
		fmt.Println("Erro ao enfileirar o lembrete:", err)
		return
	}

	_, err = s.ChannelMessageSend(reminder.ChannelID, fmt.Sprintf("Lembrete agendado para %s", timeStr))
	if err != nil {
		fmt.Println("Erro ao enviar a mensagem de confirmação: ", err)
	}
}
