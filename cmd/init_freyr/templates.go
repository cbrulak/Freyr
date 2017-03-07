package main

var surtrEnv = `SURTR_DOMAIN=https://nginx
SURTR_SECRET={{.Secret}}
SURTR_USER={{.DemoUser}}
 `

var freyrEnv = `FREYR_DEMOUSER={{.DemoUser}}
FREYR_DBHOST=postgres
FREYR_DBPASSW={{.DbPass}}
FREYR_DBUSER={{.DbUser}}
FREYR_DOMAIN={{.DomainName}}
FREYR_OAUTHID={{.OauthId}}
FREYR_OAUTHSECRET={{.OauthSecret}}
FREYR_SECRET={{.Secret}}
`

var postgresEnv = `POSTGRES_PASSWORD={{.DbPass}}
POSTGRES_USER={{.DbUser}}
`

var demoUserSQL = `insert into
   users (email, full_name, family_name, given_name, gender, locale, secret)
   values ('{{.DemoUser}}', 'demo user', 'user', 'demo', 'androgenous', 'en', '');
`

var nginxConf = `user nobody nogroup;
worker_processes auto;          # auto-detect number of logical CPU cores

events {
  worker_connections 512;       # set the max number of simultaneous connections (per worker process)
}

http {
    include mime.types;

    # thanks stackoverflow http://stackoverflow.com/a/5132440/2406040
    gzip  on;
    gzip_http_version 1.1;
    gzip_vary on;
    gzip_comp_level 6;
    gzip_proxied any;
    gzip_types text/plain text/html text/css application/json application/javascript application/x-javascript text/javascript text/xml application/xml application/rss+xml application/atom+xml application/rdf+xml;

    # make sure gzip does not lose large gzipped js or css files
    # see http://blog.leetsoft.com/2007/07/25/nginx-gzip-ssl.html
    gzip_buffers 16 8k;

    # Disable gzip for certain browsers.
    gzip_disable “MSIE [1-6].(?!.*SV1)”;

    server {
            listen         80;
            server_name    "{{.DomainName}}";
            return         301 https://$server_name$request_uri;
    }

    server {
        listen              443 ssl;
        server_name         "{{.DomainName}}";
        ssl_certificate     fullchain.cert.pem;
        ssl_certificate_key privkey.pem;
        ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
        ssl_ciphers         HIGH:!aNULL:!MD5;

        location / {
            root /usr/share/nginx/html;
            gzip_static on;
            expires 1y;
            add_header Cache-Control public;
            add_header ETag "";
            try_files $uri /index.html;
        }

        location /api/ {
            proxy_pass http://freyr:8080/;
        }
    }
}
`
