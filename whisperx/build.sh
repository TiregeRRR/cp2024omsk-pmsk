sudo docker build -t whisperx-service .
mkdir -p "$(pwd)"/whisperx-models-cache
sudo docker run -d --rm --gpus=all --env-file=.env_copy -p 8004:8000 --mount type=bind,source="$(pwd)"/whisperx-models-cache,target=/root/.cache whisperx-service