package hdr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hdr_switcher/app/internal/notify"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (

	// таймаут на запуск внешней команды
	cmdTimeout = 8 * time.Second

	// имя утилиты; если не в PATH, положите рядом с exe
	// https://github.com/res2k/HDRTray
	hdrCmdName   = "HDRCmd.exe"
	hdrCmdFolder = "HDRTray"

	StateOn          = "on"
	StateOff         = "off"
	StateUnsupported = "unsupported"
)

// Возвращает "on", "off", "unsupported".
func GetHDRState() (string, error) {
	cmdPath := ResolveHDRCmd()
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	// Предполагаем режим: HDRCmd status -m exitcode
	cmd := exec.CommandContext(ctx, cmdPath, "status", "-m", "exitcode")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // Скрывает консольное окно
	}
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()

	// Важно: нас интересует код возврата
	var exitCode int
	if err != nil {
		// Если это ошибка завершения с кодом — извлекаем
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("таймаут выполнения HDRCmd status")
		} else {
			return "", fmt.Errorf("ошибка запуска HDRCmd status: %v (%s)", err, errb.String())
		}
	} else {
		exitCode = 0
	}

	switch exitCode {
	case 0:
		return StateOn, nil
	case 1:
		return StateOff, nil
	case 2:
		return StateUnsupported, nil
	default:
		// на случай иных кодов — попробуем по тексту вывода
		txt := out.String() + errb.String()
		if containsInsensitive(txt, StateOn) {
			return StateOn, nil
		}
		if containsInsensitive(txt, StateOff) {
			return StateOff, nil
		}
		if containsInsensitive(txt, StateUnsupported) {
			return StateUnsupported, nil
		}
		return "", fmt.Errorf("неизвестный код возврата HDRCmd status: %d, вывод: %s", exitCode, txt)
	}
}

// Установить состояние HDR: true=on, false=off.
func SetHDR(on bool) error {

	cmdPath := ResolveHDRCmd()
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()

	arg := StateOff
	if on {
		arg = StateOn
	}
	cmd := exec.CommandContext(ctx, cmdPath, arg)

	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("таймаут выполнения HDRCmd %s", arg)
		}
		// если завершилось с кодом — вернём stderr
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("HDRCmd %s завершилась с кодом %d: %s", arg, ee.ExitCode(), errb.String())
		}
		return fmt.Errorf("ошибка запуска HDRCmd %s: %v (%s)", arg, err, errb.String())
	}
	return nil
}

func ResolveHDRCmd() string {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	local := filepath.Join(dir, hdrCmdFolder, hdrCmdName)
	return local
}

// Выполнить toggle: если включён — выключить, иначе включить.
func ToggleHDR() error {
	state, err := GetHDRState()
	if err != nil {
		// если статус получить не удалось — попробуем «переключить вслепую» через статус->off/on fallback
		// сначала попробуем включить
		if err2 := SetHDR(true); err2 == nil {
			notify.ShowBalloon("", "HDR включён")
			return nil
		}
		// затем выключить
		if err3 := SetHDR(false); err3 == nil {
			notify.ShowBalloon("", "HDR выключен")
			return nil
		}
		return fmt.Errorf("не удалось определить состояние и выполнить переключение: %w", err)
	}

	switch state {
	case "on":
		if err := SetHDR(false); err != nil {
			return err
		}
		notify.ShowBalloon("", "HDR выключен")
	case "off":
		if err := SetHDR(true); err != nil {
			return err
		}
		notify.ShowBalloon("", "HDR включён")
	case "unsupported":
		notify.ShowBalloon("", "HDR не поддерживается данным дисплеем/системой")
	default:
		return fmt.Errorf("неизвестное состояние: %s", state)
	}
	return nil
}

func containsInsensitive(s, sub string) bool {
	return bytes.Contains(bytes.ToLower([]byte(s)), bytes.ToLower([]byte(sub)))
}
