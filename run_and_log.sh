#!/bin/sh
LOG_DIR="/app/logs"
DEBUG_FILE="$LOG_DIR/debug_output.txt"
STRACE_LOG_FILE="$LOG_DIR/strace_output.txt"

# Ensure log directory exists (should be created by Dockerfile)
mkdir -p "$LOG_DIR"
# appuser might not have permission to chown, but Dockerfile should have set /app/logs correctly.
# Attempt to ensure files are writable by appuser if they get created by root somehow (unlikely here)
touch "$DEBUG_FILE" "$STRACE_LOG_FILE"
# chmod 666 "$DEBUG_FILE" "$STRACE_LOG_FILE" # appuser can write

echo "--- Wrapper Script Start ---" > "$DEBUG_FILE"
echo "User: $(whoami), UID: $(id -u), GID: $(id -g)" >> "$DEBUG_FILE"
echo "Executing ./nostr_bot with arg: $1 using strace" >> "$DEBUG_FILE"

strace -o "$STRACE_LOG_FILE" -f -s 256 ./nostr_bot "$1"
exit_code=$?

echo "--- strace Complete ---" >> "$DEBUG_FILE"
echo "Exit code: $exit_code" >> "$DEBUG_FILE"
echo "--- First 100 lines of strace output ($STRACE_LOG_FILE): ---" >> "$DEBUG_FILE"
head -n 100 "$STRACE_LOG_FILE" >> "$DEBUG_FILE"
echo "--- Last 100 lines of strace output ($STRACE_LOG_FILE): ---" >> "$DEBUG_FILE"
tail -n 100 "$STRACE_LOG_FILE" >> "$DEBUG_FILE"

exit $exit_code 