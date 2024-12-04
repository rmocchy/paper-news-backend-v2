.PHONY: init
init:
	cp env.example.json env.json

.PHONY: install-tools
install-tools:
	cd gin-lambda && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install github.com/xo/xo@6fe83c5

.PHONY: new-migrate
new-migrate:
	if [ -z $(name) ]; then \
		echo "name is required"; \
		exit 1; \
	fi
	sql-migrate new -config=db/dbconfig.yml $(name)

.PHONY: up-migrate
up-migrate:
	sql-migrate up -config=db/dbconfig.yml

.PHONY: gen-sql-methods
gen-sql-methods:
	cd gin-lambda && \
	rm -f ./gen/xo/* && \
	xo schema mysql://user:password@0.0.0.0:3306/paper-news_local -o ./gen/xo --go-field-tag='json:"{{ .SQLName }}" db:"{{ .SQLName }}"'
