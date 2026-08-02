# Einbuergerung-Bot: Your Ultimate (and Painless) Guide to German Citizenship! 🇩🇪🪪

🪪 **Welcome to the Einbürgerung-Bot!** Tired of dry legal texts and falling asleep over bureaucratic paperwork? Meet your witty, Telegram-powered sidekick built to help you ace the German citizenship test without losing your sanity! 🤖✨

---

## Why Einbuergerung-Bot? 🤔

- **Zero Boring Flashcards:** Level up your learning experience right inside Telegram with interactive, fun prompts. 🎲🎮
- **The Savage Bureaucracy Ranking:** Track your XP and climb the ranks - go from **Citizen in Hell** all the way to a certified, passport-wielding **German Citizen**! 📊🔥
- **Official BAMF Questions:** Master all 300 general questions plus your specific federal state (*Bundesland*) questions so you're ready for anything the exam throws at you. 🥨📜
- **100% Free & Ad-Free:** Zero hidden costs, zero premium paywalls. Just pure, unadulterated knowledge. 💸🎉

---

## Supported Commands 🕹️

### 1. `/start`
**Your Mission Begins Here!**  
Get a quick breakdown of how to conquer the citizenship test and unlock your path to that shiny new passport. 🪪🎯

### 2. `/learning`
**Let's Learn!**  
Dive headfirst into official BAMF examination questions. Earn XP, dodge the wrong answers, and climb those ranks! 📚🎮

### 3. `/progress`
**How Far Have You Come?**  
Check your current XP score and see what your current bureaucratic rank is. The finish line is in sight! 📊🔥

### 4. `/reset-federal-state`
**Did You Move?**  
Swapped your gemütlich Bavarian flat for a rainy northern coast? Use this command to switch your *Bundesland* and tackle a brand-new set of regional questions! 🗺️🔄

### 5. `/new-game`
**The Ultimate Self-Destruct Button**  
Want to wipe your progress clean at the citizenship office and start completely from scratch? Confirm the two-step prompt and dive back into the chaos! 💥🔄

### 6. Inline Callbacks
**Click Your Way to Glory!**  
Answer questions using sleek inline buttons. The bot gives you instant feedback and helpful hints when you stray off the path. 📝🕵️‍♂️

---

## For Geeks 🤓🛠️

Want to spin up the bot yourself? Let's get technical.

### 1. Starting from Sources

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/AntonSk98/Einbuergerung-Bot.git
   cd Einbuergerung-Bot
   ```

2. **Set Up Environment Variables:**
   Copy the `.env.example` file to `.env` and plug in your Telegram Bot Token.
   ```bash
   cp .env.example .env
   nano .env
   ```

3. **Install Dependencies:**
   Make sure you have Go installed, then run:
   ```bash
   go mod download
   ```

4. **Build and Run Locally:**

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

### 2. Building the Image Locally

To build the Docker image yourself, run:

```bash
docker build -t einbuergerung-bot .
```

To fire up the container manually:

```bash
docker run -e TELEGRAM_TOKEN=<your_telegram_token> -e DATABASE_PATH=/app/einbuergerung.db einbuergerung-bot
```

### 3. Using Docker Compose (The Easy Way)

In the root folder, you will find the `docker-compose.yml` file. Add your Telegram token to the environment variables, then start the application using Docker Compose:

```bash
docker-compose up -d
```

---

## FAQ ❓

**Q: Is Einbuergerung-Bot really free?** <br>
**A**: Absolutely! It's open-source and free for everyone. No corporate overlords, no sneaky paywalls! 🤑

**Q: Does it cover all topics in the test?** <br>
**A**: Yes! Every question is scraped straight from the official BAMF pool so you won't encounter any nasty surprises on exam day. 🌍📜

**Q: Will I pass the Einbürgerungstest if I train with this bot?** <br>
**A**: Definitely! 🚀🎉

---

## Getting Help 🛟

- **Found a Bug?** Open an issue on our [GitHub repository](https://github.com/AntonSk98/Einbuergerung-Bot/issues) and let us know!

## License 📜

Licensed under the Apache 2 License—check the [LICENSE](LICENSE) file for the fine print.

README.md
Displaying README.md.