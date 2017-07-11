## ng2-ui-auth-example - OAuth with [ng2-ui-auth](https://github.com/ronzeidman/ng2-ui-auth) and a Go server

### Online Demo

- I run this at [ng2-satellizer.pc-dummy.net](https://ng2-satellizer.pc-dummy.net).

### Requirements to run your own

- Go 1.7+ (I use [gvm](https://github.com/moovweb/gvm) for that).
- Nodejs (I use [nvm](https://github.com/creationix/nvm) inside a pyenv venv).
- GNU make to build the go server.
- An account on at least one of Google, Facebook, Github to test it.

### Usage

1. $ go get github.com/pcdummy/ng2-ui-auth-example
2. $ cd $GOPATH/src/github.com/pcdummy/ng2-ui-auth-example
3. $ cp secrets.ini.tmpl to secrets.ini
4. Goto Google, Facebook and Github and create an OAuth app.
5. Paste the ClientId and the Secret into your secrets.ini.
6. Install node, npm and go 1.7+
7. Run "$ make" to build the go server
8. Build the angular client: $ cd client/; npm i; npm run dev && npm run prod
9. Run the server: $ cd ..; ng2uiauthexampled --config dev.ini serve

### NGINX Config

This is my nginx.conf, in the case you want to host your own demo:

```
error_log /var/log/nginx/error.log;
events {
    use epoll;
    worker_connections 1024;
}
http {
    access_log /var/log/nginx/access.log;
    client_max_body_size 150m;
    default_type application/octet-stream;
    error_log /var/log/nginx/error.log;
    gzip on;
    gzip_static on;
    gzip_buffers 16 8k;
    gzip_comp_level 6;
    gzip_disable msie6;
    gzip_http_version 1.1;
    gzip_min_length 256;
    gzip_proxied expired no-cache no-store private auth;
    gzip_types text/plain text/css application/json application/x-javascript text/xml application/xml application/xml+rss text/javascript application/vnd.ms-fontobject application/x-font-ttf font/opentype image/svg+xml image/x-icon;
    gzip_vary on;
    keepalive_timeout 65;
    sendfile on;
    server_names_hash_bucket_size 128;
    server_tokens off;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA;
    ssl_prefer_server_ciphers on;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_session_cache shared:SSL:50m;
    ssl_session_timeout 1d;
    tcp_nodelay on;
    tcp_nopush on;
    types_hash_max_size 2048;

    include /etc/nginx/mime.types;
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
pid /run/nginx.pid;
user www-data;
worker_processes 4;
worker_rlimit_core 500M;
working_directory /tmp/;

server {
    server_name ng2-satellizer.pc-dummy.net;
    server_tokens off;
    access_log /var/log/nginx/ng2-satellizer.pc-dummy.net-access.log;
    error_log /var/log/nginx/ng2-satellizer.pc-dummy.net-error.log;
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    index index.html;
    charset utf-8;

    location ~ /.well-known {
        root /var/www/letsencrypt;
        allow all;
    }
    ssl_certificate /etc/letsencrypt/live/ng2-satellizer.pc-dummy.net/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ng2-satellizer.pc-dummy.net/privkey.pem;
    ssl_dhparam /etc/letsencrypt/live/ng2-satellizer.pc-dummy.net/dhparam.pem;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:AES:CAMELLIA:DES-CBC3-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!PSK:!aECDH:!EDH-DSS-DES-CBC3-SHA:!EDH-RSA-DES-CBC3-SHA:!KRB5-DES-CBC3-SHA;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_stapling on;
    ssl_stapling_verify on;
    add_header Strict-Transport-Security max-age=15768000;

    location / {
        index index.html;
        try_files $uri $uri/ @rewrites;
    }

    location @rewrites {
        rewrite ^ /index.html last;
    }
    root /var/www/lxdweb/static/;

    location /api/ {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_pass http://localhost:8081;
        proxy_connect_timeout 75;
        proxy_read_timeout 185;
    }
}

server {
    server_name ng2-satellizer.pc-dummy.net;
    server_tokens off;
    access_log /var/log/nginx/ng2-satellizer.pc-dummy.net-access.log;
    error_log /var/log/nginx/ng2-satellizer.pc-dummy.net-error.log;
    listen 80;
    listen [::]:80;

    location ~ /.well-known {
        root /var/www/letsencrypt;
        allow all;
    }
    return 301 https://ng2-satellizer.pc-dummy.net;
}

```


### TODO

### Authors

- René Jochum @pcdummy

### License

MIT
