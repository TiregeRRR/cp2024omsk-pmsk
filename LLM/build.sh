wget https://huggingface.co/mradermacher/T-lite-0.1-GGUF/resolve/main/T-lite-0.1.Q8_0.gguf -P ./models/
sudo docker run -d --rm --gpus=all -p 34343:8080 --name llama-cpp \
  --mount type=bind,source=./models/T-lite-instruct-0.1.Q8_0.gguf,target=/models/model-q8_0.gguf \
  --mount type=bind,source=./json_arr.gbnf,target=/grammar/json_arr.gbnf \
  --mount type=bind,source=./system_promt,target=/promt/system \
  llama.cpp:server-cuda-12.4.0 \
  -m /models/model-q8_0.gguf \
  --threads 32 \
  --batch-size 128 \
  --ubatch-size 512 \
  --ctx-size 8192 \
  --predict -1 \
  --temp 0.0 \
  --top-k 50 \
  --top-p 0.9 \
  --repeat-penalty 1.1 \
  --flash-attn \
  --mlock \
  -nkvo \
  --grammar /grammar/json_arr.gbnf \
  -ngl 99