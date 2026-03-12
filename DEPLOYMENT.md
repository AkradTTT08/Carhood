# Deployment Guide

This guide explains how to deploy the **Carhood (TTT Space Quiz)** system to a server using Docker.

## Prerequisites
- A server with **Docker** and **Docker Compose** installed.
- **Git** (if you're cloning the repository on the server).

## Project Structure
The project should have the following structure for deployment:
- `backend/` (with its `Dockerfile`)
- `frontend/` (with its `Dockerfile`)
- `docker-compose.yml` (at the root)

## Network & Ports
Ensure the following ports are open in your server's firewall:
- **Port 3000**: Frontend Access
- **Port 8081**: Backend API & WebSocket
- **Port 22**: SSH (for management)

## Deployment Steps

1. **Get the Code on the Server**
   The best way to get the code is using Git. Log into your server and run:
   ```bash
   git clone <YOUR_REPOSITORY_URL>
   cd Carhood
   ```
   *Note: If it's a private repository, you may need to [setup SSH keys](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) or use a Personal Access Token (PAT) for HTTPS.*

   Alternatively, you can upload the project folder manually using SCP or SFTP.

2. **Configure Environment (Optional)**
   The system is pre-configured to work out of the box. However, you can create a `.env` file in the root if you need to change the database URI:
   ```env
   MONGODB_URI=mongodb://mongodb:27017
   ```

3. **Build and Run**
   Navigate to the project root and run:
   ```bash
   docker-compose up -d --build
   ```
   This command will:
   - Build the Go backend.
   - Build the SvelteKit frontend (Node.js).
   - Start a MongoDB database.
   - Link everything together.

4. **Access the System**
   - **Frontend**: `http://your-server-ip:3000`
   - **Backend API**: `http://your-server-ip:8081`
   - **MonogDB**: `localhost:27017` (internal to Docker)

## Important Notes
- **Uploads Folder**: The backend automatically creates an `uploads` folder for images. This folder is persistent within the container.
- **Dynamic URLs**: The frontend is designed to automatically detect your server's IP/hostname, so you don't need to manually change URLs when moving between local and server environments.

## Troubleshooting
- To see logs: `docker-compose logs -f`
- To stop the system: `docker-compose down`
