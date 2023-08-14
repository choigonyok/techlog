FROM --platform=linux/amd64 node:16-alpine AS build

RUN mkdir app

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

FROM --platform=linux/amd64 nginx

EXPOSE 3000

COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/build /usr/share/nginx/html 