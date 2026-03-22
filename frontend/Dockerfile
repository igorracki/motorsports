FROM node:current-alpine3.23
WORKDIR /app
COPY ./ ./
RUN npm install \
    && npm run build
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1
CMD ["npm", "run", "start"]
