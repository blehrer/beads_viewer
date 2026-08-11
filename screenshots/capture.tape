# README TUI screenshots via --debug-render (reliable in CI and locally).
# Fixture copied by scripts/capture_screenshots.sh before this tape runs.

Output screenshots/.capture.gif

Set Shell bash
Set FontFamily "GeistMono Nerd Font"
Set FontSize 14
Set Width 1280
Set Height 720
Set Theme "Dracula"
Set Padding 0
Set TypingSpeed 10ms

Env BV_NO_BROWSER "1"
Env BV_NO_SAVED_CONFIG "1"
Env BV_TEST_MODE "1"
Env BV_BACKGROUND_MODE "0"
Env BV_ICON_SET "nerd"
Env TERM "xterm-256color"

Type "clear"
Enter
Type "./bv --db screenshots/fixture/.beads --theme dark --debug-render list --debug-width 120 --debug-height 40"
Enter
Sleep 1s
Screenshot screenshots/screenshot_01__main_screen.png

Type "clear"
Enter
Type "./bv --db screenshots/fixture/.beads --theme dark --debug-render insights --debug-width 120 --debug-height 40"
Enter
Sleep 1s
Screenshot screenshots/screenshot_02__insights_view.png

Type "clear"
Enter
Type "./bv --db screenshots/fixture/.beads --theme dark --debug-render board --debug-width 120 --debug-height 40"
Enter
Sleep 1s
Screenshot screenshots/screenshot_03__kanban_view.png

Type "clear"
Enter
Type "./bv --db screenshots/fixture/.beads --theme dark --debug-render graph --debug-width 120 --debug-height 40"
Enter
Sleep 1s
Screenshot screenshots/screenshot_04__graph_view.png
