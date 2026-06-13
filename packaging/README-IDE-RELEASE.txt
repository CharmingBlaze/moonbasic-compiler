moonBASIC IDE — quick start (everything in one folder)
======================================================

WHAT'S IN THIS FOLDER
---------------------

  moonbasic-ide (or moonbasic-ide.exe)  — Desktop IDE: edit .mb files, docs, check, compile, run

  moonbasic (or moonbasic.exe)          — Compiler + language server (--check, --lsp, .mb → .mbc)

  moonrun   (or moonrun.exe)            — Game runtime (F5 / Run opens your game window)

  Full moonBASIC documentation is built into the IDE (Documentation panel).

  You do NOT need Go, Node.js, VS Code, or any other tools. Extract and start coding.


FIRST STEPS
-----------

  1. Extract this zip/tar anywhere permanent (Desktop, Projects, etc.).

  2. Start the IDE:

       Windows:  double-click  START-IDE.bat
                 or run  moonbasic-ide.exe

       Linux:    chmod +x moonbasic-ide moonbasic moonrun START-IDE.sh
                 ./START-IDE.sh

       macOS:    chmod +x moonbasic-ide moonbasic moonrun START-IDE.sh
                 ./START-IDE.sh

  3. The IDE auto-detects moonbasic and moonrun in this same folder.
     Status bar should show "Toolchain ready".

  4. Write or open a .mb file, then:

       F5              Run (moonrun)
       Ctrl+Shift+C    Check syntax
       Ctrl+Shift+B    Compile to .mbc
       Alt+H           Help at cursor

  5. New project from terminal (optional):

       moonbasic new MyGame
       cd MyGame
       Open main.mb in the IDE, or:  moonrun main.mb


SETTINGS
--------

  Gear icon (title bar) or File → Settings:

    • Themes — 7 color presets + custom colors
    • Editor — font size, line height, monospace font
    • Toolchain — override compiler paths (usually leave blank)


TIPS
----

  • Begin Here: open Documentation → BEGIN_HERE.md inside the IDE.

  • Examples ship in the public moonBASIC repo: github.com/CharmingBlaze/moonbasic/tree/main/examples

  • VS Code / Cursor alternative: use moonbasic install-vscode from the full runtime zip.

  • Engine source (contributors): github.com/CharmingBlaze/moonbasic-compiler

  • User downloads & docs: github.com/CharmingBlaze/moonbasic


PLATFORM NOTES
--------------

  Windows: WebView2 is required (included on Windows 10/11). If the IDE won't start,
           install the Evergreen WebView2 Runtime from Microsoft.

  Linux:   Needs a normal desktop (GTK + WebKit). GPU drivers for games via moonrun.

  macOS:   Apple Silicon (arm64). On first launch, right-click → Open if Gatekeeper blocks.

