WORKDIR = $(shell pwd)

# 爬虫抓取网站静态页面
.PHONY: crawler
crawler:
	cd $(WORKDIR)/crawler \
	&& go mod tidy  \
	&& go run . -dist ../dist

# 编译 docker 镜像
.PHONY: image
image:
	docker build -t coolshell-nginx .

# 运行 docker
.PHONY: docker
docker:
	docker run -d --name coolshell-nginx -p 8080:80 coolshell-nginx