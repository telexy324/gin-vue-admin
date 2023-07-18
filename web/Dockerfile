FROM node:16.20.0

WORKDIR /gva_web/

COPY package.json package-lock.json .
RUN npm install

COPY . .
RUN npm run build

FROM nginx:1.24-alpine
LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"

COPY .docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf
COPY --from=0 /gva_web/dist /usr/share/nginx/html
RUN cat /etc/nginx/nginx.conf
RUN cat /etc/nginx/conf.d/my.conf
RUN ls -al /usr/share/nginx/html
