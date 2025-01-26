<img src="assets/cardArt.png" width="600" />

**All-In-Intelligence** is a text-based texas holdem poker game implemented in Go. The game leverages an interactive user interface (TUI) with visual styling.

## Description

The game is designed to be played through a console or terminal, with the option to play against an LLM (Large Language Model) player. It features an easy-to-use interface created using the TUI libraries **Bubble Tea** and **Lip Gloss**.

- **Bubble Tea**: Library used for building interactive console applications.
- **Lip Gloss**: Library used to enhance the design and styling of the user interface.

## Prerequisites

To run this project, ensure you have the following:

- Go programming language installed on your system.

## Use of the Project

1. Clone the repository:
    ```bash
    git clone https://github.com/p-obrthr/all-in-intelligence.git
    ```
2. Navigate to the source code directory:
    ```bash    
    cd all-in-intelligence/src/cmd
    ```
3. Run the project:
    ```bash
    go run .
    ```

---

## Configuration

Before starting the game, you will be prompted to configure the game. You can adjust various settings during this phase.

> Please note that you will need an OpenAI API key (e.g., as an os environment variable `OPENAI_API_KEY`) to play against an LLM player.

<img src="assets/config.png" width="300" />

## Gameplay

Once configured, the game will begin, and you will see the user interface during gameplay.

<img src="assets/gameplay.png" width="300" />