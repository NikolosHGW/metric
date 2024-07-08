// Multichecker состоящий из стандартных статических анализаторов пакета
// golang.org/x/tools/go/analysis/passes, всех анализаторов класса SA пакета staticcheck.io.
//
// Использование:
//
// В корне проекта запустите 'make my-lint'.
// Команда создаст папку my-lint с файлом result.txt.
// Для удаления папки my-lint с результатами запустить комманду 'make clear-my-lint'
package main

import (
	customAnalysis "github.com/NikolosHGW/metric/internal/analysis"
	"github.com/NikolosHGW/metric/internal/static"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	staticConfig := static.NewStaticConfig()
	staticChecks := []*analysis.Analyzer{
		// Анализатор, который обнаруживает, если в append есть только одна переменная.
		appends.Analyzer,
		// Анализатор, который сообщает о несоответствиях между файлами сборки и объявлениями Go.
		asmdecl.Analyzer,
		// Анализатор, который обнаруживает бесполезные присваивания.
		assign.Analyzer,
		// Анализатор, который проверяет распространенные ошибки при использовании пакета sync/atomic.
		atomic.Analyzer,
		// Анализатор, который проверяет аргументы функций sync/atomic на выравнивание по 64-битной границе.
		atomicalign.Analyzer,
		// Анализатор, который обнаруживает распространенные ошибки, связанные с булевыми операторами.
		bools.Analyzer,
		// Анализатор, который строит представление SSA безошибочного пакета и возвращает набор всех функций в нем.
		buildssa.Analyzer,
		// Анализатор, который проверяет теги сборки.
		buildtag.Analyzer,
		// Анализатор, который обнаруживает некоторые нарушения правил передачи указателей cgo.
		cgocall.Analyzer,
		// Анализатор, который проверяет литералы составных типов без ключей.
		composite.Analyzer,
		// Анализатор, который проверяет ошибочную передачу блокировок по значению.
		copylock.Analyzer,
		// Анализатор, который предоставляет синтаксический граф управления потоком (CFG) для тела функции.
		ctrlflow.Analyzer,
		// Анализатор, который проверяет использование reflect.DeepEqual с значениями ошибок.
		deepequalerrors.Analyzer,
		// Анализатор, который проверяет распространенные ошибки в инструкциях defer.
		defers.Analyzer,
		// Анализатор, который проверяет известные директивы инструментария Go.
		directive.Analyzer,
		// Анализатор, который проверяет, что второй аргумент errors.As является указателем на тип, реализующий error.
		errorsas.Analyzer,
		// Анализатор, который обнаруживает структуры, которые использовали бы меньше памяти, если бы их поля были отсортированы.
		fieldalignment.Analyzer,
		// Анализатор, который служит тривиальным примером и тестом API анализа.
		findcall.Analyzer,
		// Анализатор, который сообщает о коде сборки, который уничтожает указатель на фрейм до его сохранения.
		framepointer.Analyzer,
		// Анализатор, который проверяет ошибки при использовании HTTP-ответов.
		httpresponse.Analyzer,
		// Анализатор, который выделяет невозможные утверждения типа интерфейс-интерфейс.
		ifaceassert.Analyzer,
		// Анализатор, который проверяет ссылки на переменные внешнего цикла из вложенных функций.
		loopclosure.Analyzer,
		// Анализатор, который проверяет неудачу вызова функции отмены контекста.
		lostcancel.Analyzer,
		// Анализатор, который проверяет бесполезные сравнения с nil.
		nilfunc.Analyzer,
		// nilness исследует граф управления потоком функции SSA и сообщает об ошибках, таких как разыменование nil-указателя и вырожденные сравнения с nil.
		nilness.Analyzer,
		// Анализатор, который проверяет согласованность строк формата Printf и аргументов.
		printf.Analyzer,
		// Анализатор, который проверяет случайное использование == или reflect.DeepEqual для сравнения значений reflect.Value.
		reflectvaluecompare.Analyzer,
		// Анализатор, который проверяет сдвиги, превышающие ширину целого числа.
		shift.Analyzer,
		// Анализатор, который обнаруживает неправильное использование не буферизированного сигнала в качестве аргумента для signal.Notify.
		sigchanyzer.Analyzer,
		// Анализатор, который проверяет несоответствие пар ключ-значение в вызовах log/slog.
		slog.Analyzer,
		// Анализатор, который проверяет вызовы sort.Slice, которые не используют тип среза в качестве первого аргумента.
		sortslice.Analyzer,
		// Анализатор, который проверяет орфографические ошибки в подписях методов, похожих на известные интерфейсы.
		stdmethods.Analyzer,
		// Анализатор, который выделяет преобразования типов из целых чисел в строки.
		stringintconv.Analyzer,
		// Анализатор, который проверяет, что теги полей структуры правильно сформированы.
		structtag.Analyzer,
		// Анализатор для обнаружения вызовов Fatal из тестовой горутины.
		testinggoroutine.Analyzer,
		// Анализатор, который проверяет распространенные ошибочные использования тестов и примеров.
		tests.Analyzer,
		// Анализатор, который проверяет использование вызовов time.Format или time.Parse с плохим форматом.
		timeformat.Analyzer,
		// Анализатор, который проверяет передачу типов, не являющихся указателями или интерфейсами, функциям unmarshal и decode.
		unmarshal.Analyzer,
		// Анализатор, который проверяет недостижимый код.
		unreachable.Analyzer,
		// Анализатор, который проверяет недопустимые преобразования uintptr в unsafe.Pointer.
		unsafeptr.Analyzer,
		// Анализатор, который проверяет неиспользованные результаты вызовов определенных чистых функций.
		unusedresult.Analyzer,
		// Анализатор проверяет неиспользованные записи в элементы структуры или массива.
		unusedwrite.Analyzer,
		// Анализатор, который проверяет использование обобщенных функций, добавленных в Go 1.18.
		usesgenerics.Analyzer,
		// Анализатор, который проверяет, что основная функция не вызывает os.Exit().
		customAnalysis.Analyzer,
	}

	// SE анализаторы из библиотеки staticheck
	for _, check := range staticcheck.Analyzers {
		if _, ok := staticConfig.RuleSet[check.Analyzer.Name]; ok {
			staticChecks = append(staticChecks, check.Analyzer)
		}
	}

	multichecker.Main(staticChecks...)
}
