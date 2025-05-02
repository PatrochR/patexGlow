# ✨ Go TUI Text Editor (patexGlow) ✨

Hello there! Welcome to your new little editing buddy, built with Go and the amazing Charm libraries ([Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss))! 💖 It's a cozy little spot right in your terminal for editing text files.

## What Can This Little Buddy Do? 🪄

* **Jot down your thoughts:** Basic text editing fun! 📝
* **Peek at files:** A handy file browser panel shows you what's nearby. 📁
* **Open Sesame:** Pop open files from the browser with `Enter`! 🎉
* **Save Your Treasures:** Don't lose your work! `Ctrl+S` saves the day. 💾
* **Seek and Find:** Lost some text? `Ctrl+F` helps you search! 🔍
* **Magic Jumps:** See exactly where your search term is hiding (by line number) and zap! `Enter` takes you right there! ✨
* **Hop Around:** Easily switch focus between the editor, files, and search results with `Tab`. 👉
* **Handy Hints:** A little status bar tells you the filename, format, save status, and useful key hints. ℹ️
* **Pretty Colors:** Styled with love using Lipgloss for a pleasant look! 🎨

## Ingredients Needed 🧁

Just need Go (version 1.18+ is lovely!) and these awesome libraries (Go will grab them for you!):

* `github.com/charmbracelet/bubbletea`
* `github.com/charmbracelet/bubbles/textarea`
* `github.com/charmbracelet/bubbles/table`
* `github.com/charmbracelet/bubbles/textinput`
* `github.com/charmbracelet/lipgloss`

## Let's Get Set Up! 🚀

1.  **Get the code (if it's online):**
    ```bash
    git clone <your-repository-url>
    cd <your-project-directory>
    ```

2.  **Build your editor:**
    ```bash
    go build -o your_cute_editor_name .
    ```
    (Pick a fun name for `your_cute_editor_name`!)

3.  **Or just install it directly:**
    ```bash
    go install .
    ```
    (This puts it where Go keeps handy tools!)

## Time to Edit! ✏️

* **Start fresh:** Opens or creates `out.txt`.
    ```bash
    ./your_cute_editor_name
    ```
    or just:
    ```bash
    your_cute_editor_name
    ```

* **Open a specific file:**
    ```bash
    ./your_cute_editor_name path/to/your/file.txt
    ```

## Secret Handshakes (Keybindings) 🤫

* **The Essentials:**
    * `Ctrl+C`: Say goodbye! 👋
    * `Ctrl+S`: Save your masterpiece! 💾
    * `Ctrl+F`: Time to search! 🕵️‍♀️
    * `Tab`: Hop between panels! 🐇

* **In the Editor 📝:**
    * Type away! (Letters, numbers, symbols...)
    * `Arrow Keys`: Move your cursor around.

* **In the File Browser 📁:**
    * `Arrow Up/Down`: Choose a file or folder.
    * `Enter`: Open the chosen file! (Doesn't open folders... yet!)

* **In the Search Box 🔍:**
    * Type what you're looking for!
    * `Enter`: Go find it!

* **In the Search Results ✨:**
    * `Arrow Up/Down`: Pick a line number.
    * `Enter`: Zoom! Go to that line in the editor! 슝!

## Future Sparkle! (TODO / Dreams) 💭

Things we'd love to add to make this buddy even better:

* **Highlighting Magic:** Make code or search words stand out! ✨
* **Show *Where*:** Pinpoint search terms right on the line.
* **Folder Adventures:** Let's explore other directories! 🗺️
* **Make New Friends:** Create brand new files.
* **Personal Touches:** Maybe custom colors or keys?
* **Smoother Sailing:** Even better handling if things go slightly wrong.


