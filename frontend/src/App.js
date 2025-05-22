import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import ChickenManagement from './components/ChickenManagement';
import EmployeeManagement from './components/EmployeeManagement';
import Reports from './components/Reports';
import Dashboard from './components/Dashboard';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        {/* Заголовок */}
        <header className="app-header">
          <div className="container">
            <h1 className="app-title">🐔 Система управления птицефабрикой</h1>
            <nav className="nav">
              <Link to="/" className="nav-link">Главная</Link>
              <Link to="/chickens" className="nav-link">Куры</Link>
              <Link to="/employees" className="nav-link">Работники</Link>
              <Link to="/reports" className="nav-link">Отчеты</Link>
            </nav>
          </div>
        </header>

        {/* Основной контент */}
        <main className="main-content">
          <div className="container">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/chickens" element={<ChickenManagement />} />
              <Route path="/employees" element={<EmployeeManagement />} />
              <Route path="/reports" element={<Reports />} />
            </Routes>
          </div>
        </main>

        {/* Подвал */}
        <footer className="app-footer">
          <div className="container">
            <p>&copy; 2025 Система управления птицефабрикой. Курсовая работа.</p>
          </div>
        </footer>
      </div>
    </Router>
  );
}

export default App;