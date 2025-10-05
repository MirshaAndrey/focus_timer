package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

func main() {
	focusTime := 1 * time.Minute
	breakTime := 5 * time.Minute

	fmt.Println("=== Таймер фокусировки ===")
	fmt.Printf("Фокус: %v\n", focusTime)
	fmt.Printf("Отдых: %v\n", breakTime)
	fmt.Println("Нажмите Ctrl+C для выхода")
	fmt.Println("==========================")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		// Фокус
		fmt.Printf("\n⏱️  Фокус: %v\n", focusTime)
		fmt.Println("Начинаем работать...")

		focusTimer := time.NewTimer(focusTime)
		select {
		case <-focusTimer.C:
			fmt.Println("\n🔔 Время фокусировки истекло!")
			fmt.Println("🎉 Отличная работа!")
			notify("focus_end")
		case <-sigChan:
			fmt.Println("\nВыход из программы...")
			return
		}

		fmt.Printf("\n☕ Отдых: %v\n", breakTime)
		fmt.Println("Отдыхайте, отвлекитесь...")

		breakTimer := time.NewTimer(breakTime)
		select {
		case <-breakTimer.C:
			fmt.Println("\n⏰ Время отдыха истекло!")
			fmt.Println("🚀 Возвращаемся к работе!")
			notify("break_end")
		case <-sigChan:
			fmt.Println("\nВыход из программы...")
			return
		}

		fmt.Println("\n🔄 Начинаем следующий цикл...")
		time.Sleep(2 * time.Second)
	}
}

func notify(notificationType string) {
	switch notificationType {
	case "focus_end":
		fmt.Println("🔔 ФОКУС ЗАВЕРШЕН!")
		fmt.Println("🎯 Отличная работа! Сделайте пару глубоких вдохов.")
		fmt.Println("🔄 Теперь время отдыха!")
		sendSystemNotification("Фокусировка завершена", "Пора сделать перерыв!")
	case "break_end":
		fmt.Println("🔔 ОТДЫХ ЗАВЕРШЕН!")
		fmt.Println("🚀 Время вернуться к работе!")
		fmt.Println("💪 Вы готовы к новой фокусной сессии!")
		sendSystemNotification("Перерыв окончен", "Время вернуться к работе!")
	}

	printNotification()
}

func sendSystemNotification(title, message string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf(`
			Add-Type -AssemblyName System.Windows.Forms
			$balloon = New-Object System.Windows.Forms.NotifyIcon
			$balloon.Icon = [System.Drawing.SystemIcons]::Information
			$balloon.BalloonTipIcon = [System.Windows.Forms.ToolTipIcon]::Info
			$balloon.BalloonTipText = '%s'
			$balloon.BalloonTipTitle = '%s'
			$balloon.Visible = $true
			$balloon.ShowBalloonTip(10000)
		`, message, title))

	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))

	case "linux":
		_, err := exec.LookPath("notify-send")
		if err == nil {
			cmd = exec.Command("notify-send", title, message)
		} else {
			fmt.Printf("⚠️  Уведомление (Linux): %s - %s\n", title, message)
			return
		}
	default:
		fmt.Printf("⚠️  Уведомление: %s - %s\n", title, message)
		return
	}

	if cmd != nil {
		err := cmd.Run()
		if err != nil {
			fmt.Printf("⚠️  Ошибка отправки уведомления: %v\n", err)
		}
	}
}

func printNotification() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Focus Timer                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
