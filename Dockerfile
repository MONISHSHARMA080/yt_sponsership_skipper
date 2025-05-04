# —————————————————————————————————————————————————————————————
# build image for CI with non-headless Chrome support for Codespaces
# —————————————————————————————————————————————————————————————
FROM golang:1.24-bullseye

# Install Chrome, ChromeDriver, Node.js & C build tools
RUN apt-get update && apt-get install -y \
    ca-certificates curl gnupg unzip fontconfig \
    libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 \
    libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
    libxss1 libxtst6 lsb-release wget \
    build-essential \
    # X11 and display utilities for non-headless mode
    xvfb x11vnc fluxbox \
    # Additional dependencies for GUI
    libnss3 libgbm1 libasound2 \
    # Git for noVNC
    git

RUN rm -rf /var/lib/apt/lists/*

# Install noVNC for browser-based VNC access (perfect for Codespaces)
RUN git clone https://github.com/novnc/noVNC.git /opt/novnc && \
    git clone https://github.com/novnc/websockify /opt/novnc/utils/websockify && \
    ln -s /opt/novnc/vnc.html /opt/novnc/index.html

# Chrome stable
RUN wget -qO- https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" \
    > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update && apt-get install -y google-chrome-stable \
    && rm -rf /var/lib/apt/lists/*

# Create and set permissions for Chrome profile directory
RUN mkdir -p /tmp/chrome-profile && chmod 777 /tmp/chrome-profile

# ChromeDriver (match your Chrome version)
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

# Set up display for X11 forwarding
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

EXPOSE 5173 8080 9222 6080

# Create a script to start Xvfb and run tests
RUN echo '#!/bin/bash\n\
# Start Xvfb\n\
Xvfb :99 -screen 0 1280x720x16 &\n\
# Start window manager\n\
fluxbox &\n\
# Start VNC server\n\
x11vnc -display :99 -nopw -forever -shared &\n\
# Start noVNC (web-based VNC)\n\
/opt/novnc/utils/novnc_proxy --vnc localhost:5900 --listen 6080 &\n\
# Print access info\n\
echo "================================================="\n\
echo "Chrome browser UI is available via noVNC at:"\n\
echo "http://localhost:6080/"\n\
echo "================================================="\n\
# Run the original command\n\
npm run dev --prefix website -- --host 0.0.0.0 &\n\
go test ./tests -v\n\
echo " ======= SERVER LOG ======="\n\
cat tests/server.log\n\
' > /app/start.sh && chmod +x /app/start.sh

# Run the start script
CMD ["/app/start.sh"]