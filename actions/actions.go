package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func EnviarMensagem(s *discordgo.Session, m *discordgo.MessageCreate, resposta string) {
	s.ChannelMessageSend(m.ChannelID, resposta)
}

func SendHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	EnviarMensagem(s, m, "Comandos disponíveis:\n"+
		"/dailyb help - Mostra essa mensagem\n"+
		"/dailyb remember <hora> - Envia uma mensagem no chat no horário especificado\n"+
		"/dailyb status - Mostra o status do bot")
}

func ShowStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	EnviarMensagem(s, m, "Estou online e funcionando!")
}

func EnviarMensagemNoChat(s *discordgo.Session, m *discordgo.MessageCreate, trecho string) {
	if m.Author.Bot {
		EnviarMensagem(s, m, "Eu não respondo a bots! HÁ HÁ HÁ")
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
		ChannelID: m.ChannelID,
	}

	EnqueueReminder(reminder)

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lembrete agendado para %s", timeStr))
	if err != nil {
		fmt.Println("Erro ao enviar a mensagem de confirmação: ", err)
	}

	now := time.Now()
	reminderTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	duration := reminderTime.Sub(now)

	go func() {
		<-time.After(duration)
		_, err := s.ChannelMessageSend(m.ChannelID, "Hora da daily!!")
		if err != nil {
			fmt.Println("Erro ao enviar a mensagem:", err)
		}
	}()

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lembrete agendado para %s", timeStr))
	if err != nil {
		fmt.Println("Erro ao enviar a mensagem de confirmação: ", err)
	}
}
