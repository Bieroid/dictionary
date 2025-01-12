package comand

import (
	"github.com/Bieroid/dictionary/service"
	"errors"
	"fmt"
	"strings"
	"bufio"
	"os"
)

func init() {
	fmt.Printf("Добро пожаловать в словарь для перевода слов c английского.\nНа данный момент проект находится в разработке и для ознакомления предоставлена данная бета-версия(v 0.3)\nДля вывода списка поддерживаемых команд - введите help.\n\n")
}

func checkInput(input string) (err error) {
	if len(strings.Fields(input)) != 1 {
		err = errors.New("некорректный ввод")
		return
	}
	return nil
}

func ftReader(reader *bufio.Reader) (input string, err error) {
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	err = checkInput(input)
	return
}

func commandList() {
	fmt.Printf("Список комманд:\n 1. Перевод\n 2. Добавить\n 3. Помощь\n 4. Языки\n 5. Выход\n\nВведите комманду:\n")
}

func Start() {
	reader := bufio.NewReader(os.Stdin)

	for {
		commandList()
		input, err := ftReader(reader)
		fmt.Println()

		if err != nil {
			fmt.Println("Ошибка: ", err)
			continue
		}
		switch input {
			case "translate", "перевод", "1":
				handleTranslate(reader)
			case "add", "добавить", "2":
				handleAdd(reader)
			case "help", "помощь", "3":
				handleHelp()
			case "languages", "языки", "4":
				handleShowLanguages()
			case "exit", "выход", "5":
				return
			default:
				err := errors.New("некорректно введенная команда, просьба использовать команды из списка по комманде \"help\" или \"помощь\"")
				fmt.Printf("Ошибка: %s\n\n", err)
		}
	}
}

func ftReaderLang(reader *bufio.Reader) (input string, err error) {
	lang := ""

	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	err = checkInput(input)
	if err == nil {
		lang, err = service.LanguageCall(input)
		if err != nil {
			return
		}
		if lang != "" {
			input = lang
		}
	}
	return
}

func handleTranslate(reader *bufio.Reader) {
	var language string
	var word string
	var err error

	handleShowLanguages()
	fmt.Println("Введите язык:")
	language, err = ftReaderLang(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Println("Введите слово:")
	word, err = ftReader(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	trans, err := service.Translate(language, word)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Printf("Перевод слова %s: %s\n\n", word, trans)
}

func handleAdd(reader *bufio.Reader) {
	fmt.Printf("1. Для добавления языка в словарь - введите \"language\" или \"язык\"\n2. Для добавления слова в словарь - введите \"word\" или \"слово\"\n\nВведите комманду:\n")
	command, err := ftReader(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	switch command {
		case "language", "язык", "1":
			handleAddLanguage(reader)
		case "word", "слово", "2":
			handleAddWord(reader)
		default:
			err := errors.New("некорректно введенная команда, просьба использовать команды из списка выше")
			fmt.Printf("Ошибка: %s\n\n", err)
			return
	}
}

func handleAddLanguage(reader *bufio.Reader) {
	fmt.Printf("\nВведите желаемый язык для добавления в словарь:\n")
	language, err := ftReader(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	err = service.AddLanguange(language)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Printf("%s язык успешно добавлен\n\n", language)
}

func handleAddWord(reader *bufio.Reader) {
	var language string
	var word string
	var trans string
	var err error

	fmt.Println()
	handleShowLanguages()
	fmt.Printf("Введите язык:\n")
	language, err = ftReaderLang(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Println("Введите слово на русском языке:")
	word, err = ftReader(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Println("Введите перевод:")
	trans, err = ftReader(reader)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	err = service.AddWord(language, trans, word)
	if err != nil {
		fmt.Printf("Ошибка: %s\n\n", err)
		return
	}
	fmt.Printf("Слово успешно добавлено в словарь\n\n")
}

func handleShowLanguages() {
	languages := service.GetLanguages()
	fmt.Printf("Поддерживаемые языки в словаре:\n")
	for i, language := range languages {
		fmt.Printf("%d. %s\n", i + 1, language)
	}
	fmt.Println()
}

func handleHelp() {
	fmt.Printf("Список поддерживаемых комманд:\n\ntranslate - комманда для запроса перевода слова.\nadd - комманда для добавления языка или слова в словарь.\n")
	fmt.Println("languages - команда для вывода поддерживаемых языков словаря\nexit - команда для закрытия приложения\n\nВсе представленные комманды поддерживают мультиязычность и могут быть вызваны на русском языке представленным списком:")
	fmt.Printf("\ntranslate - перевод\nadd - добавить\nlanguages - языки\nexit - выход\n\nКоманды являются регистронезависимыми\nПредставленная версия словаря - 'v0.3'\n")
}
