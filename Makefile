.PHONY: tag

# make tag version=v1.0.0 description="Initial release"
tag:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make tag version=v1.0.0 description=\"Your description\""; \
		exit 1; \
	fi
	@if [ -z "$(description)" ]; then \
		echo "Error: description is required. Usage: make tag version=v1.0.0 description=\"Your description\""; \
		exit 1; \
	fi
	git tag -a $(version) -m "$(description)"
	git push origin $(version)
	@echo "Tag $(version) created and pushed successfully"