package utils

type Comandos uint8

const (
	Help Comandos = iota
	Remember
	Status
	Hi
	Flood
	Barreto
	Unknown
)

func DecodeAction(action string) Comandos {
	switch action {
	case "help":
		return Help
	case "remember":
		return Remember
	case "status":
		return Status
	case "oi":
		return Hi
	case "flood":
		return Flood
	case "barreto":
		return Barreto
	default:
		return Unknown
	}
}

func PrepareHelpMessage() string {
	message := "Olá, eu sou o Bot. Estou aqui para te ajudar a lembrar de coisas importantes.\n\n"
	message += "Para me usar, digite `d.b` seguido de um comando. Os comandos disponíveis são:\n\n"
	message += "`help` - Exibe esta mensagem de ajuda.\n"
	message += "`remember` - Cria um lembrete. Exemplo: `d.b remember 10m Lembre-se de beber água!`\n"
	message += "`status` - Exibe o status dos lembretes.\n"
	message += "`oi` - Digo olá para você.\n"
	message += "`flood` - Envia uma mensagem privada para você.\n"
	message += "`barreto` - Digo olá para o Barreto.\n"
	message += "\n"
	message += "Para mais informações, entre em contato com o Barreto."
	return message
}
