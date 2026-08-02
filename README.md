# Einbuergerung-Bot: Your Ultimate Guide to German Citizenship!

🪪 **Welcome to the Einbürgerung-Bot!** Your smart and friendly sidekick to help you ace the German citizenship test and get your passport! 🤖✨

## Why Einbuergerung-Bot?

- **Interactive Fun:** Engage with our bot through Telegram for a personalized and entertaining learning experience. No boring flashcards here! 🎲🎮
- **Gamified Progress:** Track your net XP score and level up from Citizen in Hell all the way to an official German citizen with our savage bureaucracy ranking system! 📊🔥
- **Official BAMF Questions:** Master all 300 general questions plus your specific federal state (Bundesland) questions to guarantee you're 100% prepared. 🇩🇪✨
- **100% Free:** Einbuergerung-Bot is completely free to use without any advertisements. Enjoy the learning journey! 💸🎉

## Supported Commands

### 1. `/start`
**Your Mission Begins Here!**  
Welcome to the Einbürgerung-Bot! Get a brief introduction and learn how to master the German citizenship test. 🪪🎯

### 2. `/learning`
**Let's Learn!**  
Start your personal learning mission and answer questions from the official BAMF examination libraries. Collect XP and climb through the ranks! 📚🎮

### 3. `/progress`
**How Far Have You Come?**  
Check your current XP score and find out which rank you have achieved. The path to your holy passport becomes clear! 📊🔥

### 4. `/reset-federal-state`
**Choose a New Federal State!**  
Want to change your federal state? Use this command to reset your chosen state and pick a new one. Ready for a new challenge? 🗺️🔄

### 5. `/new-game`
**Start Over!**  
Want to destroy your entire progress record at the citizenship office and start from scratch? Confirm your decision in a two-step process and embark on a new mission! 💥🔄

### 6. Inline Callbacks
**Interact with the Questions!**  
Answer the questions using the provided inline buttons. The bot gives you immediate feedback and helps you find the correct answer. Wrong answers are corrected with a short hint! 📝🕵️‍♂️

- **Answer Questions:** Choose between options A, B, C, and D to submit your answer. 🇦🇧🇨🇩
- **Select Federal State:** Choose your federal state from a list of inline buttons to receive specific questions. 🗺️📍

These interactive elements make learning with Einbürgerung-Bot not only fun but also effective! 🎲🎉

## How to Get Started

### 1. Starting from Sources

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/AntonSk98/Einbuergerung-Bot.git
   cd Einbuergerung-Bot
   ```

2. **Set Up Environment Variables:**
   Copy the `.env.example` file to `.env` and fill in your Telegram Bot Token.
   ```bash
   cp .env.example .env
   nano .env
   ```

3. **Install Dependencies:**
   Make sure you have Go installed, then run:
   ```bash
   go mod download
   ```

4. **Build and Run the Bot Locally:**

   **Linux:**
   ```bash
   go build -o binary/bot cmd/bot/main.go
   ./binary/bot
   ```

   **Windows:**
   ```bash
   go build -o binary/bot.exe cmd/bot/main.go
   binary/bot.exe
   ```

### 2. Building Image Locally

To build the Docker image locally, run the following command:

```bash
docker build -t einbuergerung-bot .
```

To start the Docker container locally, use the following command:

```bash
docker run -e TELEGRAM_TOKEN=<your_telegram_token> -e DATABASE_PATH=/app/einbuergerung.db einbuergerung-bot
```

Replace `<your_telegram_token>` with your actual Telegram Bot Token.

### 3. Using Publicly Available Image and Docker Compose

In the root folder, you will find the docker-compose.yml file. Replace `<your_telegram_token>` with your actual Telegram Bot Token in the `docker-compose.yml` file and run:

```bash
docker-compose up -d
```

## FAQ

**Q: Is Einbuergerung-Bot free to use?**
A: Absolutely! Einbuergerung-Bot is open-source and free for everyone to enjoy. No hidden fees or premium versions here! 🤑

**Q: Does Einbuergerung-Bot cover all topics in the German citizenship test?**
A: Yes, our questions are up-to-date and taken directly from BAMF to ensure you're fully prepared for the exam! 🌍📜

## Getting Help

Need help or have questions? We've got you covered!

- **Report Issues:** Found a bug or have a suggestion? Open an issue on our [GitHub repository](https://github.com/AntonSk98/Einbuergerung-Bot/issues).

## License

This project is licensed under the Apache 2 License - see the [LICENSE](LICENSE) file for details.