FROM nginx:1.27-alpine-slim

COPY public /usr/share/nginx/html
COPY nonsense /usr/share/nginx/nonsense
COPY nginx.conf /etc/nginx/nginx.conf