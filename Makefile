TARGET=jTd
PACKAGES=convert json to dart

.PHONY: all
all: build

build:
	@go build -o ./bin/$(TARGET)

clean:
	@go clean
	@rm -f ./bin/$(TARGET)

install:
	@cp ./bin/$(TARGET) /usr/local/bin

uninstall:
	@rm /usr/local/bin/$(TARGET)
