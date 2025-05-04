# ————————————————————————————————————————————————————————————— 
# build image for CI 
# —————————————————————————————————————————————————————————————
FROM golang:1.24-bullseye

# Install Chrome, ChromeDriver, Node.js & C build tools
RUN apt-get update && apt-get install -y \
    ca-certificates curl gnupg unzip fontconfig \
    libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 \
    libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
    libxss1 libxtst6 lsb-release wget \
    build-essential \
    # Install X virtual framebuffer and tools for GUI support
    xvfb x11vnc xauth xorg dbus-x11 x11-utils \
    # Add fonts for better rendering
    fonts-liberation \
    # For screenshots and better rendering
    libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 libgtk-3-0 \
    libgbm1 libasound2

RUN rm -rf /var/lib/apt/lists/*

# Chrome stable
RUN wget -qO- https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" \
    > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update && apt-get install -y google-chrome-stable \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /tmp/chrome-profile && chown -R root:root /tmp/chrome-profile

# ChromeDriver (match your Chrome version)
# Note: You should match this to your Chrome version
ARG CHROMEDRIVER_VERSION=114.0.5735.90
RUN wget -O /tmp/chromedriver.zip \
    "https://chromedriver.storage.googleapis.com/${CHROMEDRIVER_VERSION}/chromedriver_linux64.zip" \
    && unzip /tmp/chromedriver.zip -d /usr/local/bin/ \
    && chmod +x /usr/local/bin/chromedriver \
    && rm /tmp/chromedriver.zip

# Node.js 23.x
RUN curl -fsSL https://deb.nodesource.com/setup_23.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Teach cgo to link against libdl
ENV CGO_ENABLED=1 
ENV CGO_LDFLAGS="-ldl"

# Set display for X server
ENV DISPLAY=:99

# cache Go deps
COPY go.mod go.sum ./
RUN go mod download

# cache website deps
COPY website/package*.json website/
RUN npm ci --prefix website

# global tools
RUN npm install -g wait-on localtunnel

# copy everything (including .env)
COPY . .

# Create a wrapper script to start Xvfb before running tests
RUN echo '#!/bin/bash\n\
# Start Xvfb\n\
Xvfb :99 -screen 0 1920x1080x24 -ac &\n\
XVFB_PID=$!\n\
\n\
# Give Xvfb time to start\n\
sleep 1\n\
\n\
# Optional: Start VNC server to view the browser (if needed)\n\
# x11vnc -display :99 -forever -nopw &\n\
\n\
# Run the original command\n\
npm run dev --prefix website -- --host 0.0.0.0 &\n\
\n\
# Run tests\n\
go test ./tests -v\n\
\n\
# Show logs\n\
echo " ======= SERVER LOG ======="\n\
cat tests/server.log\n\
\n\
# Kill Xvfb\n\
kill $XVFB_PID\n\
' > /app/start.sh && chmod +x /app/start.sh

EXPOSE 5173 8080 
# Optional: Expose VNC port if you want to view the browser remotely
EXPOSE 5900

# Run both front-end and Go tests using our wrapper script
CMD ["/app/start.sh"]