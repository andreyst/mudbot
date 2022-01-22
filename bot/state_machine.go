package bot

// TODO: Add scheduling event loop iteration to make it independent of data from mud
func (b *Bot) ProcessEvent(e Event) {
	if e == EVENT_NOP {
		return
	}

	if b.Fight.IsActive {
		return
	}

	b.logger.Debugf("Processing event %v", e)

	switch b.State {
	case STATE_INITIALIZING:
		switch e {
		case EVENT_ENCODING_MENU:
			b.SendToMud("u")
		case EVENT_LOGIN_PROMPT:
			b.SendToMud(b.Credentials.Login)
		case EVENT_PASSWORD_PROMPT:
			b.SendToMudWithoutEcho(b.Credentials.Password)
		case EVENT_PRESS_CRFL_PROMPT:
			b.SendToMud("")
		case EVENT_GAME_MENU:
			b.SendToMud("1")
		case EVENT_PROMPT:
			b.SendToMud("score")
			b.SwitchState(STATE_IDLE)
		}
	case STATE_IDLE:
		b.ProcessIdle()
	case STATE_DRINKING:
		switch e {
		case EVENT_DRANK:
			if b.Char.IsThirsty {
				b.Drink()
			} else {
				b.SwitchState(STATE_IDLE)
			}
		case EVENT_WATER_CONTAINER_EMPTY:
			b.ErrorClient("Water container empty!")
		case EVENT_NO_SUCH_ITEM:
			b.SwitchState(STATE_STUCK)
		}
	case STATE_EATING:
		switch e {
		case EVENT_ATE:
			if b.Char.IsHungry {
				b.Eat()
			} else {
				b.SwitchState(STATE_IDLE)
			}
		case EVENT_NO_SUCH_ITEM:
			b.SwitchState(STATE_STUCK)
		}
	case STATE_RESTING:
		switch e {
		case EVENT_TICK:
			if b.Char.Stamina == b.Char.MaxStamina || b.Char.IsHungry || b.Char.IsThirsty {
				b.Stand()
				b.SwitchState(STATE_IDLE)
			}
		}
	}
}

func (b *Bot) ProcessEvents(events []Event) {
	for _, e := range events {
		b.ProcessEvent(e)
	}
}

func (b *Bot) ProcessIdle() {
	if b.Char.IsThirsty {
		b.SwitchState(STATE_DRINKING)
		b.Drink()
	} else if b.Char.IsHungry {
		b.SwitchState(STATE_EATING)
		b.Eat()
	} else if b.Char.Stamina < 15 {
		b.SwitchState(STATE_RESTING)
		b.Rest()
	}
}

func (b *Bot) SwitchState(newState State) {
	b.logger.Debugf("state: %v->%v", b.State, newState)
	switch newState {
	case STATE_IDLE:
		b.SendToMud("") // CR;LF to trigger prompt
	case STATE_STUCK:
		b.ErrorClient("I AM STUCK")
	}
	b.State = newState
}
