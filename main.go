package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"dailybot/actions"
	"dailybot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("KEY"))
	if err != nil {
		fmt.Println("Erro ao criar a sessão do Discordgo:", err)
		return
	}

	dg.AddHandler(onReady)
	dg.AddHandler(onMessage)

	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir a sessão do Discordgo:", err)
		return
	}

	fmt.Println("Bot está online. Aguardando comandos...")

	// Aguardar por interrupção do programa (Ctrl+C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	dg.Close()
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Bot está pronto! Logado como", r.User.Username)
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorar mensagens próprias do bot
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Message.Content, "d.b") {
		executarAcao(s, m, m.Message.Content)
	}
}

func executarAcao(s *discordgo.Session, m *discordgo.MessageCreate, trecho string) {
	trecho = strings.Replace(trecho, "d.b", "", 1)
	trecho = strings.TrimSpace(trecho)

	action := utils.DecodeAction(trecho)
	switch action {
	case utils.Help:
		actions.SendHelp(s, m)
	case utils.Remember:
		actions.EnviarMensagemNoChat1(s, m, trecho)
		go actions.ProcessReminders(s)
	case utils.Status:
		actions.ShowStatus(s, m)
	case utils.Hi:
		actions.EnviarMensagem(s, m, "Olá, como posso ajudar?")
	case utils.Flood:
		actions.SendMessageUser(s, m)
	case utils.Barreto:
		actions.EnviarMensagem(s, m, "Barreto, para de barretice")
	default:
		actions.EnviarMensagem(s, m, "Não entendi.. Poderia repetir? (X_X)")
	}
}
