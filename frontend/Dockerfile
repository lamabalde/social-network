FROM node:lts-alpine

WORKDIR /frontend

# Copy only the necessary files initially (package.json and package-lock.json)
COPY package*.json ./

# Install dependencies
RUN npm install

COPY . .

RUN npm run build

EXPOSE 8080

CMD ["npm", "run", "dev"]
#CMD ["npm", "start"]