# CSMS Dashboard

A web dashboard for the Charging Station Management System (CSMS).
This dashboard provides real-time monitoring, management, and analytics for EV charging stations, users, transactions, and more.

---

## 🚀 Features

- **Charge Point Monitoring**: View status, connectors, and live updates for all charge points.
- **Transaction Management**: Track, filter, and analyze charging sessions and costs.
- **User & RFID Management**: Manage users and their RFID tags.
- **Role-based Access**: Admin and operator views.
- **Responsive UI**: Works on desktop and mobile.
- **Authentication**: Secure login with JWT.
- **OCPP Integration**: Real-time updates via backend OCPP events.

---

## 🛠️ Tech Stack

- **Frontend Framework**: Vue 3 + Composition API
- **Build Tool**: Vite
- **UI Library**: Vuetify 3
- **State Management**: Pinia
- **Routing**: Vue Router
- **HTTP Client**: Axios
- **Linting**: ESLint, Pretti

---

## 📦 Project Structure

```
csms-dashboard/
├── public/                # Static assets
├── src/
│   ├── assets/            # Images, icons, etc.
│   ├── components/        # Vue components (tables, forms, charts, etc.)
│   ├── pages/             # Page-level views (Dashboard, Transactions, etc.)
│   ├── router/            # Vue Router config
│   ├── services/          # API service modules (Axios)
│   ├── stores/            # Pinia stores
│   ├── utils/             # Utility functions
│   ├── App.vue            # Root component
│   └── main.js            # Entry point
├── package.json
├── vite.config.mjs
└── README.md
```

---

## ⚡ Quick Start

1. **Install dependencies**
   ```bash
   npm install
   # or
   yarn install
   ```

2. **Configure environment**
   - Edit `src/services/api.js` or `.env` if you use environment variables for API base URL.

3. **Run the development server**
   ```bash
   npm run dev
   # or
   yarn dev
   ```

4. **Build for production**
   ```bash
   npm run build
   # or
   yarn build
   ```

---

## 🔌 Backend Integration

- This dashboard is designed to work with the [CSMS Backend](../README.md) (OCPP 1.6, REST API).
- Make sure the backend server is running and accessible at the configured API URL.

---

## 🧪 Testing & Linting

- **Lint code:**  
  `npm run lint`
- **Format code:**  
  `npm run format`

---

## 📄 License

MIT License

---