package quiz

import (
	"encoding/binary"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"math/rand/v2"
)

var (
	MinComplexity = flag.Int("min-complexity", 1, "minimal complexity of generated expressions")
	MaxComplexity = flag.Int("max-complexity", 7, "maximum complexity of generated expressions")
)

var manualQuestions = []Question{
	{Text: "window.reload() = ?", Answers: []Answer{
		{Text: "undefined"},
		{Text: "null"},
		{Text: "Не равен", IsCorrect: true},
	}},
	{Text: `"hello".substring(3) = ?`, Answers: []Answer{
		{Text: "world"},
		{Text: "lo"},
		{Text: `"lo"`, IsCorrect: true},
	}},
	{Text: "На каком языке вы предпочитаете здороваться?", Answers: []Answer{
		{Text: "Haskell", IsCorrect: true},
		{Text: "Python", IsCorrect: true},
		{Text: "Французкий", IsCorrect: true},
	}},
	{Text: "Какой принтер вы бы выбрали для печати наклеек?", Answers: []Answer{
		{Text: "Струйный цветной", IsCorrect: true},
		{Text: "Термотрансферный", IsCorrect: true},
		{Text: "3D-принтер"},
	}},
	{Text: "rm -rf / = ?", Answers: []Answer{
		{Text: "--no-preserve-root", IsCorrect: true},
		{Text: "Выдернуть шнур питания", IsCorrect: true},
		{Text: "Enter-hit", IsCorrect: true},
	}},
	{Text: "Как вы думаете, могут ли технологии помочь в решении мировых проблем, связанных с религией?", Answers: []Answer{
		{Text: "Да", IsCorrect: true},
		{Text: "Нет", IsCorrect: false},
		{Text: "42", IsCorrect: false},
	}},
	{Text: "Как вы относитесь к идее создания искусственной религии с помощью технологий?", Answers: []Answer{
		{Text: "Кстати, приглашаю вступить в мою...", IsCorrect: true},
		{Text: "Не отношусь"},
		{Text: "Сегодня Тихие Филины", IsCorrect: true},
	}},
	{Text: "А завтра наступит?", Answers: []Answer{
		{Text: "Да, и оно будет прекрасным!", IsCorrect: true},
		{Text: "Да, и оно будет хуже, чем сегодня", IsCorrect: true},
		{Text: "Нет, сегодня последний день", IsCorrect: true},
	}},
	{Text: "За секунду до какого события вы бы не хотели оказаться?", Answers: []Answer{
		{Text: "Падение Римской империи", IsCorrect: true},
		{Text: "Начало Первой мировой войны", IsCorrect: true},
		{Text: "X", IsCorrect: true},
	}},
	{Text: "Как вы относитесь к роботам?", Answers: []Answer{
		{Text: "Они классные", IsCorrect: true},
		{Text: "Они меня пугают", IsCorrect: true},
		{Text: "Они лучше людей", IsCorrect: true},
	}},
	{Text: "Какой ваш любимый праздник?", Answers: []Answer{
		{Text: "День сурка", IsCorrect: true},
		{Text: "День рождения", IsCorrect: true},
		{Text: "UgraCTF Open", IsCorrect: true},
	}},
	{Text: "Какая у вас главная мечта?", Answers: []Answer{
		{Text: "Захватить мир", IsCorrect: true},
		{Text: "Быть Вечно молодым", IsCorrect: true},
		{Text: "Найти алмаз", IsCorrect: true},
	}},
	{Text: "Какой ваш любимый способ решения CTF-задач?", Answers: []Answer{
		{Text: "Брутфорс", IsCorrect: true},
		{Text: "Социальная инженерия", IsCorrect: true},
		{Text: "Криптоанализ", IsCorrect: true},
	}},
	{Text: "Какой метод сбора данных вы считаете наиболее эффективным?", Answers: []Answer{
		{Text: "Анкеты и опросы", IsCorrect: true},
		{Text: "Интервью", IsCorrect: true},
		{Text: "Анализ поведения пользователей", IsCorrect: true},
	}},
	{Text: "Какое приветствие вам нравится больше и заинтересовывает вас?", Answers: []Answer{
		{Text: "Позвольте рассказать вам об уцуцуге", IsCorrect: true},
		{Text: "Что вы знаете о CTF?", IsCorrect: true},
		{Text: "Вы решали задачи от Калана?", IsCorrect: true},
	}},
	{Text: "Как вы относитесь к идее работы в условиях высокой неопределённости?", Answers: []Answer{
		{Text: "Мне это нравится", IsCorrect: true},
		{Text: "Мне это не нравится"},
		{Text: "Мне все равно", IsCorrect: true},
	}},
	{Text: "Как вы предпочитаете готовиться к выполнению сложных задач?", Answers: []Answer{
		{Text: "Немедленное выполнение", IsCorrect: true},
		{Text: "Замедленное выполнение", IsCorrect: true},
		{Text: "Медленное выполнение", IsCorrect: true},
	}},
	{Text: "Какой подход к разработке ПО вы считаете наиболее продуктивным?", Answers: []Answer{
		{Text: "Agile", IsCorrect: true},
		{Text: "Waterfall", IsCorrect: true},
		{Text: "CTF-Solving", IsCorrect: true},
	}},
	{Text: "Какие характеристики делает клининг идеальным для вашей компании?", Answers: []Answer{
		{Text: "Пунктуальность", IsCorrect: true},
		{Text: "Тщательность", IsCorrect: true},
		{Text: "Доброжелательность", IsCorrect: true},
	}},
	{Text: "Какие личные качества вы цените в своих коллегах?", Answers: []Answer{
		{Text: "Честность", IsCorrect: true},
		{Text: "Ответственность", IsCorrect: true},
		{Text: "Отсчитываются об успехах в чате"},
	}},
	{Text: "Как вы относитесь к шуткам про программистов?", Answers: []Answer{
		{Text: "Шутил, шучу и буду шутить из шкафа", IsCorrect: true},
		{Text: "Не слушаю", IsCorrect: true},
		{Text: "Я монитор перверну", IsCorrect: true},
	}},
	{Text: "Товарищ начальник, ...", Answers: []Answer{
		{Text: "привет!", IsCorrect: true},
		{Text: "у нас тут тепло!", IsCorrect: true},
		{Text: "пюре на обед!", IsCorrect: true},
	}},
	{Text: "Какой совет вы бы дали себе завтра, исходя из сегодняшнего опыта?", Answers: []Answer{
		{Text: "Будь терпеливым", IsCorrect: true},
		{Text: "Планируй время", IsCorrect: true},
		{Text: "Учись на ошибках", IsCorrect: true},
	}},
}

func generateQuestion(seed int) Question {
	var seedBytes [32]byte
	copy(seedBytes[:], binary.LittleEndian.AppendUint64(nil, uint64(seed)))
	copy(seedBytes[8:], binary.LittleEndian.AppendUint64(nil, uint64(seed)))
	copy(seedBytes[16:], binary.LittleEndian.AppendUint64(nil, uint64(seed)))
	copy(seedBytes[24:], binary.LittleEndian.AppendUint64(nil, uint64(seed)))

	r := rand.New(rand.NewChaCha8(seedBytes))

	if r.IntN(10) < 3 {
		q := manualQuestions[r.IntN(len(manualQuestions))]
		answers := make([]Answer, len(q.Answers))
		copy(answers, q.Answers)
		r.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})
		return Question{
			Text:    q.Text,
			Answers: answers,
		}
	}

	complexity := r.IntN(*MaxComplexity-*MinComplexity) + *MinComplexity
	sb := &strings.Builder{}
	answer := generateIntExpr(r, complexity, sb)
	answers := []Answer{
		{Text: strconv.Itoa(answer), IsCorrect: true},
		{Text: strconv.Itoa(answer - 1 - r.IntN(10))},
		{Text: strconv.Itoa(answer + 1 + r.IntN(10))},
	}
	r.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})
	return Question{
		Text:    sb.String() + " = ?",
		Answers: answers,
	}
}

func generateIntExpr(r *rand.Rand, complexity int, sb *strings.Builder) int {
	op := r.IntN(4)

	if complexity < 1 || op == 0 {
		n := r.IntN(2000) - 1000
		if n < 0 {
			fmt.Fprintf(sb, "(%d)", n)
		} else {
			fmt.Fprint(sb, n)
		}
		return n
	}

	leftComplexity := complexity - r.IntN(complexity+1)
	rightComplexity := complexity - leftComplexity
	sb.WriteByte('(')
	left := generateIntExpr(r, leftComplexity, sb)
	switch op {
	case 1:
		sb.WriteByte('+')
	case 2:
		sb.WriteByte('-')
	case 3:
		sb.WriteByte('*')
	}
	right := generateIntExpr(r, rightComplexity, sb)
	sb.WriteByte(')')
	switch op {
	case 1:
		return left + right
	case 2:
		return left - right
	case 3:
		return left * right
	default:
		panic("unreachable")
	}
}
