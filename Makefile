TARGET := app
all:
	go build -o $(TARGET) cmd/bingo_bot/main.go 
	./$(TARGET)

clean:
	rm -rf $(TARGET)
	rm -rf logs/app.log