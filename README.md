# Einbuergerung-Bot: Your Ultimate Guide to German Citizenship!

🪪 **Welcome to the Einbürgerung-Bot!** Your smart and friendly sidekick to help you ace the German citizenship test and get your passport! 🤖✨

## Why Einbuergerung-Bot?

- **Interactive Fun:** Engage with our bot through Telegram for a personalized and entertaining learning experience. No boring flashcards here! 🎲🎮
- **Gamified Progress:** Track your net XP score and level up from Citizen in Hell all the way to an official German citizen with our savage bureaucracy ranking system! 📊🔥
- **Official BAMF Questions:** Master all 300 general questions plus your specific federal state (Bundesland) questions to guarantee you're 100% prepared. 🇩🇪✨
- **100% Free:** Einbuergerung-Bot is completely free to use without any advertisements. Enjoy the learning journey! 💸🎉

## How to Get Started

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
   Build and start the bot with the following commands:
   ```bash
   go build -o binary/bot cmd/bot/main.go
   ./binary/bot
   ```

   **Windows:**
   Build and start the bot with the following commands:
   ```bash
   go build -o binary/bot.exe cmd/bot/main.go
   binary/bot.exe
   ```

5. **Interact with the Bot:**
   Open Telegram and search for your bot. Start chatting and learning with Einbuergerung-Bot!

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
