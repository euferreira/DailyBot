package actions

import "github.com/bwmarrin/discordgo"

func EnviarMensagem(session *discordgo.Session, m *discordgo.MessageCreate, text string) {
	_, error := session.ChannelMessageSend(m.ChannelID, text)
	if error != nil {
		println("Erro ao enviar a mensagem: ", error)
		panic("Erro ao enviar a mensagem" + error.Error())
	}
}
