# —————————————————————————————————————————————————————————————
# build image for CI
# (Located at ./tests/Dockerfile)
# —————————————————————————————————————————————————————————————
FROM golang:1.24-bullseye

# install Chrome, ChromeDriver, Node.js & C build tools
# No changes needed in this section
RUN apt-get update && apt-get install -y \
      ca-certificates curl gnupg unzip fontconfig \
      libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 \
      libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
      libxss1 libxtst6 lsb-release wget \
      build-essential \
    && rm -rf /var/lib/apt/lists/*

# Chrome stable
# No changes needed in this section
RUN wget -qO- https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
  && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" \
         > /etc/apt/sources.list.d/google-chrome.list \
  && apt-get update && apt-get install -y google-chrome-stable \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /tmp/chrome-profile && chown -R root:root /tmp/chrome-profile

# ChromeDriver (match your Chrome version)
# Consider updating this version if your Chrome version changes
ARG CHROMEDRIVER_VERSION=114.0.5735.90
RUN wget -O /tmp/chromedriver.zip \
         "https://chromedriver.storage.googleapis.com/${CHROMEDRIVER_VERSION}/chromedriver_linux64.zip" \
  && unzip /tmp/chromedriver.zip -d /usr/local/bin/ \
  && chmod +x /usr/local/bin/chromedriver \
  && rm /tmp/chromedriver.zip

# Node.js 23.x (Consider using an LTS version like 20.x or 22.x for stability)
RUN curl -fsSL https://deb.nodesource.com/setup_23.x | bash - \
  && apt-get install -y nodejs \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Teach cgo to link against libdl
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-ldl"

# cache Go deps
# These paths are relative to the build context (project root), which is correct.
COPY go.mod go.sum ./
RUN go mod download

# cache website deps
# This path is relative to the build context (project root), which is correct.
COPY website/package*.json website/
RUN npm ci --prefix website

# global tools
RUN npm install -g wait-on localtunnel

# copy everything (including .env if present)
# Copies from the build context (project root) to /app inside the container. Correct.
COPY . .

EXPOSE 5173 8080

# run both front-end and Go tests
# Paths here are relative to WORKDIR /app inside the container. Correct.
CMD sh -c "\
    npm run dev --prefix website -- --host 0.0.0.0 & \
    go test ./tests -v ;\
    echo '  ======= SERVER LOG =======' && \
    cat tests/server.log  \
"
