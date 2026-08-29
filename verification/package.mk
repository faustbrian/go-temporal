.PHONY: docs php-compat

docs:
	./scripts/check-docs.sh

php-compat:
	test -n "$(PHP_TEMPORAL_SOURCE)" || { echo "set PHP_TEMPORAL_SOURCE"; exit 2; }
	generated="$$(mktemp)"; trap 'rm -f "$$generated"' EXIT; \
	php compat/php/generate.php "$(PHP_TEMPORAL_SOURCE)" > "$$generated"; \
	php compat/php/compare.php "$$generated" compat/fixtures/php_v1.json
