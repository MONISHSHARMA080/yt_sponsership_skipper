# —————————————————————————————————————————————————————————————
# build image for CI
# —————————————————————————————————————————————————————————————
FROM golang:1.24-bullseye

# install Chrome, ChromeDriver, Node.js & C build tools
RUN apt-get update && apt-get install -y \
      ca-certificates curl gnupg unzip fontconfig \
      libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 \
      libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
      libxss1 libxtst6 lsb-release wget \
      build-essential 
     
RUN rm -rf /var/lib/apt/lists/*
 

# Chrome stable
RUN wget -qO- https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
  && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" \
       > /etc/apt/sources.list.d/google-chrome.list \
  && apt-get update && apt-get install -y google-chrome-stable \
  && rm -rf /var/lib/apt/lists/*

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

EXPOSE 5173 8080

# run both front-end and Go tests
CMD sh -c "\
    npm run dev --prefix website -- --host 0.0.0.0 & \
    go test ./tests -v \
"
