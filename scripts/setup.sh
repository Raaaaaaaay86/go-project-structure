work_dir=$(pwd | sed 's/\//\\\//g')

# Change to current directory
sed -i '' "s/\${WORK_DIR}/$work_dir/" config/app_config.yaml

# If brew not exists then install it
if ! which brew >/dev/null 2>&1; then
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# If ffmpeg not exists then install it
if ! which ffmpeg >/dev/null 2>&1; then
  brew install ffmpeg
fi

docker-compose up -d