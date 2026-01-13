<div id="top"></div>

<!-- PROJECT LOGO -->
<br />
<div align="center">
    <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" alt="Logo" width="80" height="80">

  <h3 align="center">Inference API - LLM & TTS Service</h3>

  <p align="center">Real-time League of Legends game event monitor that connects to the Live 
Client API, tracks game events, stores them in PostgreSQL, and forwards them 
to an inference service for AI processing.</p>
</div>

 <br />  

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">🧭 About The Project</a>
      <ul>
        <li><a href="#built-with">🏗️ Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">📋 Getting Started</a>
      <ul>
        <li><a href="#prerequisites">🗺️ Prerequisites</a></li>
        <li><a href="#installation">⚙️ Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">💾 Usage</a></li>
    <li><a href="#api-endpoints">🔌 API Endpoints</a></li>
    <li><a href="#architecture">🏛️ Architecture</a></li>
    <li><a href="#contributing">🔗 Contributing</a></li>
    <li><a href="#license">📰 License</a></li>
    <li><a href="#contact">📫 Contact</a></li>
    <li><a href="#acknowledgments">⛱️ Acknowledgments</a></li>
  </ol>
</details>

<br>

<!-- ABOUT THE PROJECT -->
## 🧭 About The Project

This project is a real-time event monitoring service for League of Legends 
games. It connects to the Riot Games Live Client API (running locally during 
a game), captures game events, stores them in a database, and forwards them 
to an inference API for AI-powered commentary generation.

KEY FEATURES:
- **Real-time Event Monitoring**: Connects to League of Legends Live Client API
- **Event Tracking**: Monitors game events (kills, objectives, turrets, etc.)
- **Database Storage**: PostgreSQL integration for event persistence
- **Player Filtering**: Watch specific players or all players
- **Event Forwarding**: Sends events to inference API for processing
- **Concurrent Processing**: Uses goroutines for efficient event handling
- **Cross-platform**: Builds for Windows, macOS (Intel & Apple Silicon)

### 🏗️ Built With

* [![PostgreSQL Badge](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
* [![DigitalOcean Badge](https://img.shields.io/badge/DigitalOcean-0080FF?style=for-the-badge&logo=DigitalOcean&logoColor=white)](https://www.digitalocean.com/)
* [![Golang Badge](https://img.shields.io/badge/Golang-0080FF?style=for-the-badge&logo=Golang&logoColor=blue)](https://golang.org/)

<p align="right"><a href="#top">⬆️</a></p>


<!-- GETTING STARTED -->
## 📋 Getting Started

To get a local copy up and running follow these simple example steps.

### 🗺️ Prerequisites

* [mise](https://mise.jdx.dev/)
* Go 1.24.4 or higher
* PostgreSQL database
* League of Legends game client (for Live Client API)
* Access to inference API endpoint (droplet) see https://github.com/Ouilliam-21/inference

### ⚙️ Installation

1. Clone the repository
```sh
git clone https://github.com/Ouilliam-21/lol-event-broker.git
cd lol-event-broker
```

2. **Install dependencies using mise**
```sh
mise trust
mise install
```

3. Install dependencies
```sh
mise exec -- go mod download
```

4. Having a postgreSQL database running
See https://github.com/Ouilliam-21/database.git

5. Configure the application
Copy the example config:
```sh
cp config.example.yaml config.yaml
```

Edit config.yaml with your settings (see Configuration section below).

1. Set up environment variables
Create a .env file if you prefer environment variables:
```sh
   DATABASE_HOST=localhost
   DATABASE_PORT=5432
   DATABASE_NAME=lol_events
   DATABASE_USER=your_user
   DATABASE_PASSWORD=your_password
```
   
<p align="right"><a href="#top">⬆️</a></p>


<!-- USAGE EXAMPLES -->
## 💾 Usage

### Development Mode

Run the application locally:
```sh
mise run run:local
#Or directly with Go
go run ./cmd/main.go -config=config.yaml
```

### Production Mode

To having a production executable run mise command followed by your pc architecture


Ex: macOS (Apple Silicon):
```sh
mise run build:mac-arm
# Run binary
./build/lol-event-mac-arm64 -config=config.yaml
```

PREREQUISITES FOR RUNNING:

1. League of Legends game - The Live Client API is only available 
   when a game is active, the project will wait until game start
2. Ensure the inference API is running - The droplet endpoint must be 
   accessible
3. Database must be running - PostgreSQL connection is required
   

<p align="right"><a href="#top">⬆️</a></p>

## 🔌 Configuration

The application uses a YAML configuration file. Example config.yaml:
```yaml
endpoints:
  live_client: "https://127.0.0.1:2999"  # League of Legends Live Client API
  droplet: "http://localhost:8000"        # Inference API endpoint
  auth_token: "your_bearer_token"         # Authentication token for droplet

database:
  host: "localhost"
  port: "5432"
  name: "lol_events"
  user: "your_username"
  password: "your_password"

events:
  # Events to watch for
  watch:
    - ChampionKill
    - BaronKill
    - DragonKill
    - HeraldKill
    - TurretKilled
    - InhibKilled
    - Ace
    - MultiKill
    - FirstBrick

players:
  # Specific players to watch (empty = all players)
  watch:
    - "Player name"
```

<p align="right"><a href="#top">⬆️</a></p>

## 🏛️ Architecture

The application follows a clean architecture pattern:
```sh
lol-event/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database repositories
│   ├── riot/                # Riot API client
│   │   ├── liveclient.go    # Live Client API connection
│   │   └── events/          # Event handling, factory
│   ├── droplet.go           # Inference API client
│   └── utils/               # Utility functions
├── config.yaml              # Configuration file
└── go.mod                   # Go dependencies
```

Data flow:

1. LiveClient: Connects to League of Legends Live Client API, polls for events
2. EventManager: Processes and filters game events
3. Database Repositories: Store game sessions and events in PostgreSQL
4. Droplet Client: Forwards event IDs to inference API
5. Goroutines: Concurrent processing of events and API calls

```sh
League of Legends Game
    ↓
Live Client API (127.0.0.1:2999)
    ↓
LoL Event Monitor
    ↓
    ├──→ PostgreSQL Database (storage)
    └──→ Inference API (processing)
```
<p align="right"><a href="#top">⬆️</a></p>

## 🔗 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right"><a href="#top">⬆️</a></p>


## 📰 License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right"><a href="#top">⬆️</a></p>


## 📫 Contact

Reach me at : gauron.dorian.pro@gmail.com

Project Link: [https://github.com/Ouilliam-21/lol-event-broker.git](https://github.com/Ouilliam-21/lol-event-broker.git)

<p align="right"><a href="#top">⬆️</a></p>


## ⛱️ Acknowledgments

This space is a list to resources i found helpful and would like to give credit to.

* [Riot Games Live Client API](https://developer.riotgames.com/docs/lol)
* [Go Documentation](https://go.dev/doc/)
* [pgx - PostgreSQL Driver](https://github.com/jackc/pgx)
* [mise - Runtime Version Manager](https://mise.jdx.dev/)

<p align="right"><a href="#top">⬆️</a></p>

<a href="https://github.com/othneildrew/Best-README-Template">Template inspired by othneildrew</a>
