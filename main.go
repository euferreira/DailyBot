package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"dailybot/actions"

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
	if strings.HasPrefix(m.Message.Content, "/dailyb") {
		executarAcao(s, m, m.Message.Content)
	}
}

func executarAcao(s *discordgo.Session, m *discordgo.MessageCreate, trecho string) {
	if strings.Contains(trecho, "/dailyb remember") {
		actions.EnviarMensagemNoChat1(s, m, trecho)
		go actions.ProcessReminders(s)
	} else if strings.Contains(trecho, "/dailyb help") {
		actions.SendHelp(s, m)
	} else if strings.Contains(trecho, "/dailyb oi") {
		actions.EnviarMensagem(s, m, "Olá, como posso ajudar?")
	} else if strings.Contains(trecho, "/dailyb status") {
		actions.ShowStatus(s, m)
	} else if strings.Contains(trecho, "/dailyb barreto") {
		actions.EnviarMensagem(s, m, "Barreto, vai curtir tuas férias! XD")
	} else {
		actions.EnviarMensagem(s, m, "Não entendi.. Poderia repetir? (X_X)")
	}
}
