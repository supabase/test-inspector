FROM node:lts as builder

WORKDIR /app

COPY . .

RUN npm install \
  --prefer-offline \
  --frozen-lockfile \
  --non-interactive \
  --production=false

RUN npm run build

RUN rm -rf node_modules && \
  NODE_ENV=production npm install \
  --prefer-offline \
  --pure-lockfile \
  --non-interactive \
  --production=true

FROM gcr.io/distroless/nodejs:16

WORKDIR /app

COPY --from=builder /app  .

ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

CMD [ ".output/server/index.mjs" ]