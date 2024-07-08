lint_dir = ./my-lint
lint_result_file = result.txt
lint_exec = mylint

.PHONY: my-lint
my-lint: _run-lint _clear_lint_binary

.PHONY: _clear_lint_binary
_clear_lint_binary:
	rm $(lint_exec)

.PHONY: _run-lint
_run-lint: _create-lint-dir _build_linter
	-./$(lint_exec) ./... 2> $(lint_dir)/$(lint_result_file)

.PHONY: _build_linter
_build_linter:
	go build -o $(lint_exec) cmd/staticlint/main.go

.PHONY: _create-lint-dir
_create-lint-dir:
	mkdir -p $(lint_dir)

.PHONY: clear-my-lint
clear-my-lint:
	rm -rf $(lint_dir)
