docker build . -t registry.cn-hangzhou.aliyuncs.com/fadinglight/bill-go:dev

docker push registry.cn-hangzhou.aliyuncs.com/fadinglight/bill-go:dev

ssh fadinglight "cd /root/docker/bill-sys/;
                docker compose down; 
                docker pull registry.cn-hangzhou.aliyuncs.com/fadinglight/bill-go:dev; 
                docker compose up -d"
