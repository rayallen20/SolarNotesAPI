# BookStack 部署指南(Ubuntu 22.04 + Nginx 源码部署)

> 在局域网服务器上以**源码方式**部署 BookStack 文档平台,并使用 Nginx 对外提供访问。
> 用于替代 showDoc 作为 API 文档平台。

**最后更新:2026-06-15**

---

## 一、部署环境

| 项目 | 说明 |
| --- | --- |
| 服务器 IP | `192.168.1.60`(局域网) |
| 操作系统 | Ubuntu 22.04 Server |
| Web 服务器 | Nginx |
| PHP | 8.3(通过 `ondrej/php` PPA 安装) |
| 数据库 | MySQL 8.0 |
| BookStack | `release` 分支(基于 Laravel 12) |
| 访问地址 | `http://192.168.1.60:4061` |

---

## 二、环境说明(与 showDoc 的关键差异)

BookStack 是基于 Laravel 的应用,部署方式与 showDoc 有以下几点不同,务必注意:

- **要求 PHP ≥ 8.2**(源码 `composer.json` 限定 `^8.2.0`)。Ubuntu 22.04 自带的是 PHP 8.1,**必须通过 `ondrej/php` PPA 升级**,本文使用 8.3。
- **依赖不随仓库打包**,需要用 **Composer** 安装(`composer install`)。请使用**官方 Composer 安装器**,**不要用 Ubuntu 的 apt `composer` 包**——它依赖 `/usr/share/php` 系统库,安装依赖时会因缺 `intl` 等扩展直接 fatal。
- **需要 MySQL / MariaDB 数据库**(showDoc 默认使用开箱即用的 SQLite)。
- **Nginx 的 web 根目录指向项目下的 `public/` 子目录**,而非项目根目录。
- 前端资源在 `release` 分支中已编译好,**无需安装 Node/npm**。

> 官方提供一键安装脚本,但其使用 Apache。本文采用 Nginx 手动部署,更可控且便于统一管理。

---

## 三、部署步骤

### 1. 安装 Nginx + PHP 8.3 + MySQL + 相关扩展

```bash
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository -y ppa:ondrej/php
sudo apt update

sudo apt install -y nginx mysql-server git unzip \
  php8.3-fpm php8.3-cli php8.3-curl php8.3-mbstring php8.3-xml \
  php8.3-zip php8.3-gd php8.3-mysql

sudo systemctl enable --now php8.3-fpm nginx

# 确认 PHP-FPM 的 socket 名,后面 Nginx 配置要用到,通常是 /run/php/php8.3-fpm.sock
ls /run/php/*.sock
```

> 扩展清单对应 `composer.json` 要求的 curl / dom / fileinfo / gd / json / mbstring / xml / zip
> (其中 dom 在 `php8.3-xml` 内,fileinfo、json 为内置),外加连接数据库用的 `php8.3-mysql`。

### 2. 创建 MySQL 数据库和用户

Ubuntu 22.04 装好的 MySQL 8 默认使用 socket 认证,`sudo mysql` 可免密进入:

```bash
sudo mysql <<'SQL'
CREATE DATABASE bookstack CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'bookstack'@'localhost' IDENTIFIED WITH mysql_native_password BY '换成你的强密码';
GRANT ALL PRIVILEGES ON bookstack.* TO 'bookstack'@'localhost';
FLUSH PRIVILEGES;
SQL
```

### 3. 拉取源码并安装依赖

```bash
cd /var/www
sudo git clone https://github.com/BookStackApp/BookStack.git --branch release --single-branch bookstack
cd /var/www/bookstack

# 准备 .env(先复制,下一步再填数据库信息)
sudo cp .env.example .env
```

接着安装**官方 Composer**(⚠️ 不要用 Ubuntu 的 apt `composer` 包:它依赖 `/usr/share/php` 系统库,安装依赖时会因缺 `intl` 等扩展而 fatal):

```bash
cd /tmp
php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
# 校验安装器完整性(官方推荐步骤)
php -r "if (hash_file('sha384','composer-setup.php')===trim(file_get_contents('https://composer.github.io/installer.sig'))){echo 'OK'.PHP_EOL;}else{unlink('composer-setup.php');echo 'CORRUPT'.PHP_EOL;exit(1);}"
sudo php composer-setup.php --quiet --install-dir=/usr/local/bin --filename=composer
rm -f composer-setup.php
composer --version
```

最后安装 BookStack 的 PHP 依赖(以 root 运行需要此环境变量以免告警):

```bash
cd /var/www/bookstack
sudo COMPOSER_ALLOW_SUPERUSER=1 composer install --no-dev
```

### 4. 配置 .env 并初始化应用

编辑 `/var/www/bookstack/.env`,修改以下几项:

```ini
APP_URL=http://192.168.1.60:4061

DB_HOST=127.0.0.1
DB_DATABASE=bookstack
DB_USERNAME=bookstack
DB_PASSWORD=刚才设置的强密码
```

> ⚠️ `APP_URL` **必须带上端口 `:4061`**:BookStack 用它拼接资源/链接地址,端口对不上会导致样式丢失、链接错乱。

生成应用密钥并执行数据库迁移(建表):

```bash
cd /var/www/bookstack
sudo php artisan key:generate --force
sudo php artisan migrate --force
```

### 5. 设置目录权限

PHP-FPM 以 `www-data` 用户运行,将代码交给它,并放开几个需要写入的目录:

```bash
sudo chown -R www-data:www-data /var/www/bookstack
sudo chmod -R 755 /var/www/bookstack
sudo chmod -R 775 /var/www/bookstack/storage \
                  /var/www/bookstack/bootstrap/cache \
                  /var/www/bookstack/public/uploads
sudo chmod 640 /var/www/bookstack/.env     # 含密码与密钥,收紧权限
```

### 6. 配置 Nginx

新建 `/etc/nginx/sites-available/bookstack`,注意 `root` 指向 **public** 子目录:

```nginx
server {
    listen 4061;  # 访问端口;改此处需同步改步骤四的 APP_URL
    server_name 192.168.1.60;          # 有域名就换成域名

    root  /var/www/bookstack/public;   # ← 指向 public,不是项目根
    index index.php index.html;

    client_max_body_size 100M;         # 允许上传较大图片/附件,按需调大

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/run/php/php8.3-fpm.sock;   # 与步骤 1 查到的 socket 一致
    }

    # 安全:禁止访问隐藏文件(放行 .well-known 以便日后签发 HTTPS 证书)
    location ~ /\.(?!well-known).* {
        deny all;
    }
}
```

启用并重载:

```bash
sudo ln -s /etc/nginx/sites-available/bookstack /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

### 7. 配置防火墙(如启用了 ufw)

```bash
sudo ufw allow 22/tcp      # 别把自己的 SSH 关在外面
sudo ufw allow 4061/tcp
sudo ufw reload
```

### 8. 首次访问与登录

浏览器打开 **http://192.168.1.60:4061**,BookStack 默认管理员账号:

- 邮箱:`admin@admin.com`
- 密码:`password`

> ⚠️ **登录后第一件事**:点击右上角头像 → Edit Profile,修改邮箱和密码。
> 随后可在 Settings 中将界面语言设为简体中文、修改站点名称等。

---

## 四、运维与可选配置

### 公网访问

`192.168.1.60` 是局域网地址,完成上述步骤后,**同网段机器**访问 `http://192.168.1.60:4061` 即可使用。
若需**公网访问**(外网也能打开),还需二选一:

- 在路由器上做 80 端口映射(前提是宽带有公网 IP);
- 使用 frp / 花生壳等内网穿透工具。

### 大附件上传

如需上传超过 100M 的文件,需同步调大两处限制:

1. Nginx 的 `client_max_body_size`;
2. PHP 的 `upload_max_filesize` 和 `post_max_size`(`/etc/php/8.3/fpm/php.ini`)。

修改后重启:`sudo systemctl restart php8.3-fpm`。

### 版本升级

BookStack 基于 Laravel,升级流程:

```bash
cd /var/www/bookstack
sudo -u www-data git pull origin release
sudo COMPOSER_ALLOW_SUPERUSER=1 composer install --no-dev
sudo php artisan migrate --force
```

> 官方提供 `php bin/run-updates.sh` 一键升级脚本,可直接使用。

### 数据备份

需备份两部分:

- **数据库**:`mysqldump -u bookstack -p bookstack > bookstack-$(date +%F).sql`
- **上传文件**:`/var/www/bookstack/storage` 和 `/var/www/bookstack/public/uploads`

### 启用 HTTPS(可选)

绑定域名后即可自动签发并配置 Let's Encrypt 证书:

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx
```

---

## 五、安全提醒

- 默认账号 `admin@admin.com` / `password`,**部署后立即修改**。
- `.env` 含应用密钥与数据库密码,权限保持 `640`,不要提交到代码仓库。
- 数据库用户使用强口令。
- 如无需公网访问,仅在局域网内开放 4061 端口即可。

---

## 参考来源

- BookStack 官方仓库:<https://github.com/BookStackApp/BookStack>
- 手动安装文档:<https://www.bookstackapp.com/docs/admin/installation/>
- PHP 版本与扩展要求来源:源码 `composer.json`(`"php": "^8.2.0"`)
