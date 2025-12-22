#!/usr/bin/env sh
set -e

# Create non-login system user 'glance'
if ! id -u glance >/dev/null 2>&1; then
  # Prefer useradd if available (Debian/Ubuntu/RHEL)
  if command -v useradd >/dev/null 2>&1; then
    NOLOGIN_SHELL="/usr/sbin/nologin"
    [ -x "$NOLOGIN_SHELL" ] || NOLOGIN_SHELL="/sbin/nologin"
    useradd --system --no-create-home --shell "$NOLOGIN_SHELL" glance || true
  # Fallback to adduser (Alpine)
  elif command -v adduser >/dev/null 2>&1; then
    NOLOGIN_SHELL="/sbin/nologin"
    [ -x "$NOLOGIN_SHELL" ] || NOLOGIN_SHELL="/usr/sbin/nologin"
    # -S (system), -H (no home), -D (no password), -s (shell)
    adduser -S -H -D -s "$NOLOGIN_SHELL" glance || true
  fi
fi

# Ensure group exists
if ! getent group glance >/dev/null 2>&1; then
  if command -v groupadd >/dev/null 2>&1; then
    groupadd --system glance || true
  elif command -v addgroup >/dev/null 2>&1; then
    addgroup -S glance || true
  fi
fi

# Make sure the glance user is in the glance group
if command -v usermod >/dev/null 2>&1; then
  usermod -a -G glance glance 2>/dev/null || true
fi

# Systemd daemon reload so units are recognized
if command -v systemctl >/dev/null 2>&1; then
  systemctl daemon-reload || true
fi

# Keep services disabled by default; just inform the user
echo "Glance systemd units installed. To enable:"
echo "  systemctl enable glance-agent.service"
echo "Then start with:"
echo "  systemctl start glance-agent.service"
echo "An example configuration file has been placed at /usr/lib/glance-agent/config.env.example"
echo "Copy it to /etc/glance-agent/config.env and modify as needed."
echo "The service will not start until a valid configuration file is in place."