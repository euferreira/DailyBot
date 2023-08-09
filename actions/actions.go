package actions

import (
	"fmt"
	"strings"
	"time"

	"dailybot/utils"

	"github.com/bwmarrin/discordgo"
)

func SendHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := utils.PrepareHelpMessage()

	EnviarMensagem(s, m, message)
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

func getUserID(s *discordgo.Session, guildID, username string) (string, error) {
	members, err := s.GuildMembers(guildID, "", 500)
	if err != nil {
		return "", err
	}
	for _, m := range members {
		fmt.Println(m.User.Username + " " + m.User.Discriminator)
		if (m.User.Username + "#" + m.User.Discriminator) == username {
			return m.User.ID, nil
		}

		if m.User.Username == username || m.Nick == username {
			return m.User.ID, nil
		}
	}
	return "", fmt.Errorf("couldn't find user %s", username)
}

func SendMessageUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Obtém o ID do canal em que a mensagem foi recebida
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Erro ao obter informações do canal:", err)
		return
	}

	// Obtém as informações do servidor
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Erro ao obter informações do servidor:", err)
		return
	}

	// Obtém o nome do servidor a partir das informações obtidas
	guildName := guild.Name
	print(guildName)

	userID, err := getUserID(s, channel.GuildID, "José Barreto#3586")
	if err != nil {
		EnviarMensagem(s, m, "Não te achei, @José Barreto")
		fmt.Println("Erro ao pegar o ID do usuário:", err)
		return
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		EnviarMensagem(s, m, "Oi <@"+userID+">!")
	}
}
