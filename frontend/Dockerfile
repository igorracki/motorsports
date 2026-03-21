FROM node:current-alpine3.23
WORKDIR /app
COPY ./ ./
RUN npm install \
    && npm run build
EXPOSE 3000
CMD ["npm", "run", "start"]
