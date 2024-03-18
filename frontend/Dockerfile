FROM node:18.15-alpine3.17 as base
WORKDIR /frontend
COPY package*.json ./
RUN npm ci \
    && npm cache clean --force
FROM base as dev
ENV NODE_ENV=development
ENV PATH /frontend/node_modules/.bin:$PATH
RUN npm i --only=development \
    && npm cache clean --force
RUN npm config list
COPY . .
CMD ["npm", "start"]