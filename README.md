# Discord Clone App

Discord clone built to replicate the functionality of the popular communication platform. This application includes all the main features you'd expect from Discord, including account creation, server and channel management, and real-time communication.

## Key Features

- **Account Management**: Create and manage user accounts.
- **Server Management**: Create, join, and leave servers.
- **Text Channels**: Create and participate in text-based chat channels.
- **Voice Channels**: Join voice channels to talk with others in real-time.
- **Video Channels**: Join video meetings for face-to-face communication.
- **Roles and Permissions**: Manage roles such as Admin, Moderator, and Guest with customizable permissions.
- **Real-Time Communication**: Seamless text, voice, and video interactions.

## Technologies Used

### Frontend
- **[Next.js](https://nextjs.org/)**: React framework for server-side rendering and building the frontend.
- **[Tailwind CSS](https://tailwindcss.com/)**: Utility-first CSS framework for styling.
- **[ShadCN](https://ui.shadcn.com/)**: Component library to enhance UI/UX design.

### Backend
- **[GoLang](https://golang.org/)**: Backend programming language.
- **[Gin](https://gin-gonic.com/)**: Web framework for building REST APIs.
- **[Gorm](https://gorm.io/)**: ORM library for Go to interact with the database.
- **[WebSocket](https://pkg.go.dev/github.com/gorilla/websocket)**: Real-time, bi-directional communication for text, voice, and video channels.
- **[WebRTC](https://pion.ly/)**: Enables real-time peer-to-peer audio, video, and data communication.

#### WebRTC Implementation
- The backend includes a **Selective Forwarding Unit (SFU)**, which optimizes WebRTC communications.
- The SFU allows the server to receive media streams from multiple clients and selectively forward them to other participants based on various conditions, such as the client’s subscription or bandwidth capabilities.
- This architecture is more scalable and efficient than peer-to-peer mesh networks, especially for video channels and large group calls.

### Database
- **[PostgreSQL](https://www.postgresql.org/)**: Relational database for storing app data.

### Deployment
- **[Docker Compose](https://docs.docker.com/compose/)**: Simplifies deployment and running of multi-container Docker applications.

## Demo

Check out the live demo of the application here:  
[Live Demo](https://www.jkrn.me/)

## Website Sample
![Website Screenshot](https://github.com/MajorTom3K1M/discord-clone/blob/main/screenshot/screenshot-1.png)

## Getting Started

### Run with Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/MajorTom3K1M/discord-clone.git
   cd discord-clone
   ```

2. Edit the environment variables in the `docker-compose.yml` file as needed.

3. Start the application using Docker Compose:
   ```bash
   docker-compose up
   ```

### Run Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/MajorTom3K1M/discord-clone.git
   cd discord-clone
   ```

2. Install PostgreSQL:
   - Follow the [official PostgreSQL installation guide](https://www.postgresql.org/download/) for your operating system.
   - Create a new database and user with the required credentials.
   - Update the `.env` file in the `backend` folder with your database credentials.

3. Configure Environment Variables:
   - The project uses two `.env` files: one in the `backend` folder and one in the `frontend` folder.
   - Each folder contains a `.env.example` file that provides basic values for the required variables. Copy and rename it to `.env` in each folder:
     ```bash
     cd backend
     cp .env.example .env
     # Edit the .env file with your PostgreSQL credentials and other necessary values

     cd ../frontend
     cp .env.example .env
     # Edit the .env file with your UploadThing credentials and other required values
     ```

4. UploadThing Integration:
   - The frontend requires `UPLOADTHING_SECRET` and `UPLOADTHING_APP_ID` to use UploadThing services.
   - Obtain these credentials from your [UploadThing dashboard](https://uploadthing.com/) and add them to the `.env` file in the `frontend` folder:
     ```
     UPLOADTHING_SECRET=your-uploadthing-secret
     UPLOADTHING_APP_ID=your-uploadthing-app-id
     ```

5. Install the dependencies for both backend and frontend:

   From the root of the project:
   ```bash
   cd frontend
   npm install
   ```

   Then, from the root of the project:
   ```bash
   cd backend
   go mod tidy
   ```

6. Run the backend and frontend concurrently:

   From the root of the project:
   ```bash
   npm install
   npm run dev
   ```

Ensure your PostgreSQL database is running and accessible on the configured host and port.
