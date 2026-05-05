FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .

ARG VITE_BASE_URL=http://localhost:8080
ARG VITE_WS_BASE_URL=ws://localhost:8080
ARG VITE_BASE_URL_HOME_ENDPOINT=home
ARG VITE_CLOSE_ENDPOINT=close
ARG VITE_WS_BASE_URL_ACTIONS_ENDPOINT=ws-actions
ARG VITE_SITE_TITLE=Splitzy
ARG VITE_ABOUT_ME_URI=

ENV VITE_BASE_URL=$VITE_BASE_URL \
    VITE_WS_BASE_URL=$VITE_WS_BASE_URL \
    VITE_BASE_URL_HOME_ENDPOINT=$VITE_BASE_URL_HOME_ENDPOINT \
    VITE_CLOSE_ENDPOINT=$VITE_CLOSE_ENDPOINT \
    VITE_WS_BASE_URL_ACTIONS_ENDPOINT=$VITE_WS_BASE_URL_ACTIONS_ENDPOINT \
    VITE_SITE_TITLE=$VITE_SITE_TITLE \
    VITE_ABOUT_ME_URI=$VITE_ABOUT_ME_URI

RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
RUN printf 'server {\n    listen 80;\n    root /usr/share/nginx/html;\n    index index.html;\n    location / {\n        try_files $uri $uri/ /index.html;\n    }\n}\n' > /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
