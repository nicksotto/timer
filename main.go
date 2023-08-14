package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
	"github.com/saaste/pomodoro/timer"
)

type TimerState int

const (
	TimerStateWait TimerState = iota
	TimerStateWork
	TimerStateShortBreak
	TimerStateLongBreak
)

type model struct {
	timerState       TimerState
	shortBreakCount  int
	pomodoroCount    int
	notificationSent bool
	keys             keyMap
	help             help.Model
	timer            timer.Timer
	progress         progress.Model
	config           Config
}

func (m model) Init() tea.Cmd {
	return timer.TickEverySecond()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		m.timer.Update()
		m.SendNotificationIfNeeded()
		return m, timer.TickEverySecond()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.StartWork):
			if m.timerState != TimerStateWork {
				m.timerState = TimerStateWork
				m.pomodoroCount += 1
				m.keys.StartWork.SetEnabled(false)
				m.keys.StartBreak.SetEnabled(true)
				m.timer.Reset(m.config.WorkDuration())
				m.notificationSent = false
			}
		case key.Matches(msg, m.keys.StartBreak):
			m.keys.StartWork.SetEnabled(true)
			m.keys.StartBreak.SetEnabled(false)
			if m.timerState == TimerStateWork {
				if m.shortBreakCount < m.config.ShortBreaks {
					m.timerState = TimerStateShortBreak
					m.shortBreakCount += 1
					m.timer.Reset(m.config.ShortBreakDuration())
				} else {
					m.timerState = TimerStateLongBreak
					m.shortBreakCount = 0
					m.timer.Reset(m.config.LongBreakDuration())
				}
				m.notificationSent = false
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("%s\n", titleStyle.Render("POMODORO TIMER"))

	switch m.timerState {
	case TimerStateShortBreak:
		if m.timer.IsTimeout() {
			s += alertStyle.Render("Time to go back to work!")
		} else {
			s += fmt.Sprintf("%s\n\n", normalStyle.Render("Time for a quick break!"))
			s += fmt.Sprintf("%s: %s\n", hilightStyle.Render("Time left"), m.timer.TimeLeft())
			s += fmt.Sprintf("%s: %s", hilightStyle.Render("Progress"), m.progress.ViewAs(m.timer.PercentageLeft()))

		}
	case TimerStateLongBreak:
		if m.timer.IsTimeout() {
			s += alertStyle.Render("Time to go back to work!")
		} else {
			s += fmt.Sprintf("%s\n\n", normalStyle.Render("Time for a longer break! You've earned it!"))
			s += fmt.Sprintf("%s: %s\n", hilightStyle.Render("Time left"), m.timer.TimeLeft())
			s += fmt.Sprintf("%s: %s", hilightStyle.Render("Progress"), m.progress.ViewAs(m.timer.PercentageLeft()))

		}
	case TimerStateWork:
		if m.timer.IsTimeout() {
			s += alertStyle.Render("Good job! Time to have a break!")
		} else {
			var nth = "1st"
			if m.pomodoroCount == 2 {
				nth = "2nd"
			} else if m.pomodoroCount > 2 {
				nth = fmt.Sprintf("%dth", m.pomodoroCount)
			}

			text := fmt.Sprintf("This is your %s pomodoro. Time to focus on the task!", nth)
			s += fmt.Sprintf("%s\n\n", normalStyle.Render(text))
			s += fmt.Sprintf("%s: %s\n", hilightStyle.Render("Time left"), m.timer.TimeLeft())
			s += fmt.Sprintf("%s: %s", hilightStyle.Render("Progress"), m.progress.ViewAs(m.timer.PercentageLeft()))
		}
	case TimerStateWait:
		s += fmt.Sprintf("%s\n\n", normalStyle.Render("When you are ready, press w to start your first pomodoro."))
	}

	s += fmt.Sprintf("\n\n%s", m.help.View(m.keys))
	return s
}

func (m *model) SendNotificationIfNeeded() {
	if !m.config.SendNotification {
		return
	}

	if m.timerState != TimerStateWait && m.timer.IsTimeout() && !m.notificationSent {
		switch m.timerState {
		case TimerStateWork:
			if m.shortBreakCount < 4 {
				beeep.Notify("Pomodoro Timer", "Time to have a quick break!", "assets/pomodoro.png")
			} else {
				beeep.Notify("Pomodoro Timer", "Time to have a longer break!", "assets/pomodoro.png")
			}
		case TimerStateShortBreak, TimerStateLongBreak:
			beeep.Notify("Pomodoro Timer", "Time to go back to work!", "assets/pomodoro.png")
		}

		m.notificationSent = true
	}
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		log.Fatalf("unable to read the config file: %s", err)
	}

	model := model{
		timer:      timer.New(config.WorkDuration()),
		timerState: TimerStateWait,
		progress:   progress.New(progress.WithScaledGradient(gradientColors[0], gradientColors[1])),
		keys:       keys,
		help:       help.New(),
		config:     config,
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run the program: %s", err)
	}
}
