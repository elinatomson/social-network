FROM node:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package*.json ./

# Install project dependencies and remove cache to reduce image size
RUN npm install --only=production && npm cache clean --force

# Copy the rest of the application code
COPY . .

# Build the React application
RUN npm run build

# Expose the port the application runs on 
EXPOSE 3000

# Start the application
CMD ["npm", "start"]
