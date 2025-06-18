<h1 align="center", style="font-size: 50px;">Movie Rating API</h1>

<p align="center">
  <img src="https://img.shields.io/badge/GO-1.22.4-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/PostgreSQL-16-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-Mit-green?style=for-the-badge" />
</p>
<br><br>

# Background
This project is made for the purpose of learning and practicing **Go language**, and the backend server side development. Its implementation is inspired by the user requirements you can find in <a href="https://github.com/MahmoodAhmed-SE/movie-rating-api/blob/main/user-requirements.md">user-requirements</a>. However, I will not follow the requirements strictly.

### Note:
<mark>For learning and testing purposes, free and open license films/videos will be used.</mark>

<br><br>

## 📌 Featuring

---

### 👤 User Management
- **POST `/api/v1/register-user`** — Register a new user.
- **POST `/api/v1/login-user`** — Authenticate and log in a user.

---

### 🎬 Movie Management
- **GET `/api/v1/movies`** — Retrieve a list of all available movies.
- **POST `/api/v1/movies`** — Post a movie.
- **GET `/api/v1/movies/{movieId}`** — Get detailed information about a specific movie.
- **GET `/api/v1/movies-rating/{movieId}`** — Fetch the average rating of a specific movie.
- **GET `/api/v1/chat-on-movie/{movieId}`** — View chat messages related to a specific movie.
- **GET `/api/v1/search/{movieName}`** — 🔧 *[Planned]* Search for a movie by name.¹

---

### 🤝 User Interaction
- **POST `/api/v1/movies-rating`** — Submit a movie rating.
- **POST `/api/v1/chat-on-movie`** — Post a message in a movie's chat.
- **POST `/api/v1/watchlist`** — Add a movie to the user's watchlist.
- **GET `/api/v1/watchlist`** — Retrieve the authenticated user's watchlist.

---
🔧 *[planned]*
### 📊 Analytics & Recommendations
- **GET `/api/v1/recommend-movies`** — Get personalized movie recommendations.
- **GET `/api/v1/get-viewer-stats`** — Retrieve statistics on viewer behavior and engagement.
