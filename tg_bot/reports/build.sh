sudo docker build -t reports .
mkdir -p "$(pwd)"/reports
sudo docker run -d --rm  -p 8005:8000 reports