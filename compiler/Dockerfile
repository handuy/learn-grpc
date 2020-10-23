FROM node:13-alpine
WORKDIR /app
COPY node.sh ./
RUN touch index.js && chmod 755 index.js
ENTRYPOINT ["./node.sh"]