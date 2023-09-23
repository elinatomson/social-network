# Use the official Node.js image as a parent image
FROM node:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React application
RUN npm run build

# Expose the port the application runs on (usually 3000 by default)
EXPOSE 3000

# Start the application
CMD ["npm", "start"]
