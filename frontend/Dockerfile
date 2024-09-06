# Stage 1: Build (in a Linux environment)
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code and build
COPY . .
RUN npm run build

# Stage 2: Production (in a Linux environment)
FROM node:18-alpine
WORKDIR /app

# Install only production dependencies
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# Copy build artifacts from the builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

# Expose the application port
EXPOSE 3000
ENV NODE_ENV production

# Start the application
CMD ["npm", "run", "start"]
