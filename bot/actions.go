package bot

func (b *Bot) Drink() {
	b.SendToMud("пить мех")
}

func (b *Bot) Eat() {
	b.SendToMud("есть хлеб")
}

func (b *Bot) Stand() {
	b.SendToMud("вста")
}

func (b *Bot) Rest() {
	b.SendToMud("отд")
}
