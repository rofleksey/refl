.PHONY: all clean

ANTLR = antlr4
GRAMMAR = parser/Refl.g4

all: gen

gen: $(GRAMMAR)
	@echo "Generating parser from $(GRAMMAR)..."
	$(ANTLR) -Dlanguage=Go -o parser/gen -package gen -visitor -no-listener $(GRAMMAR)
	@mv parser/gen/parser/* parser/gen/
	@rm -rf parser/gen/parser
	@echo "Parser generated in parser/gen/ directory"

clean:
	rm -rf parser/gen