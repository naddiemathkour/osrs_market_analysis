FROM node:22-alpine

WORKDIR /app

COPY . .

# Angular RUN Commands
RUN npm install -g @angular/cli
RUN npm install

# Expose port
EXPOSE 4200

# Host Angular server
CMD ["ng", "serve", "--host", "0.0.0.0" ]