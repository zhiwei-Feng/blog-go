docker rm vueblog
docker build -t vueblog .
docker run -p 8080:8080 --name vueblog vueblog
